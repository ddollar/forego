package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const defaultPort = 5000
const defaultShutdownGraceTime = 3

var flagPort int
var flagConcurrency string
var flagRestart bool
var flagShutdownGraceTime int
var envs envFiles

var cmdStart = &Command{
	Run:   runStart,
	Usage: "start [process name] [-f procfile] [-e env] [-p port] [-c concurrency] [-r] [-t shutdown_grace_time]",
	Short: "Start the application",
	Long: `
Start the application specified by a Procfile. The directory containing the
Procfile is used as the working directory.

The following options are available:

  -f procfile  Set the Procfile. Defaults to './Procfile'.

  -e env       Add an environment file, containing variables in 'KEY=value', or
               'export KEY=value', form. These variables will be set in the
               environment of each process. If no environment files are
               specified, a file called .env is used if it exists.

  -p port      Sets the base port number; each process will have a PORT variable
               in its environment set to a unique value based on this. This may
               also be set via a PORT variable in the environment, or in an
               environment file, and otherwise defaults to 5000.

  -c concurrency
               Start a specific number of instances of each process. The
               argument should be in the format 'foo=1,bar=2,baz=0'. Use the
               name 'all' to set the default number of instances. By default,
               one instance of each process is started.

  -r           Restart a process which exits. Without this, if a process exits,
               forego will kill all other processes and exit.

  -t shutdown_grace_time
               Set the shutdown grace time that each process is given after
               being asked to stop. Once this grace time expires, the process is
               forcibly terminated. By default, it is 3 seconds.

If there is a file named .forego in the current directory, it will be read in
the same way as an environment file, and the values of variables procfile, port,
concurrency, and shutdown_grace_time used to change the corresponding default
values.

Examples:

  # start every process
  forego start

  # start only the web process
  forego start web

  # start every process specified in Procfile.test, with the environment specified in .env.test
  forego start -f Procfile.test -e .env.test

  # start every process, with a timeout of 30 seconds
  forego start -t 30
`,
}

func init() {
	cmdStart.Flag.StringVar(&flagProcfile, "f", "Procfile", "procfile")
	cmdStart.Flag.Var(&envs, "e", "env")
	cmdStart.Flag.IntVar(&flagPort, "p", defaultPort, "port")
	cmdStart.Flag.StringVar(&flagConcurrency, "c", "", "concurrency")
	cmdStart.Flag.BoolVar(&flagRestart, "r", false, "restart")
	cmdStart.Flag.IntVar(&flagShutdownGraceTime, "t", defaultShutdownGraceTime, "shutdown grace time")
	err := readConfigFile(".forego", &flagProcfile, &flagPort, &flagConcurrency, &flagShutdownGraceTime)
	handleError(err)
}

func readConfigFile(config_path string, flagProcfile *string, flagPort *int, flagConcurrency *string, flagShutdownGraceTime *int) error {
	config, err := ReadConfig(config_path)

	if config["procfile"] != "" {
		*flagProcfile = config["procfile"]
	} else {
		*flagProcfile = "Procfile"
	}
	if config["port"] != "" {
		*flagPort, err = strconv.Atoi(config["port"])
	} else {
		*flagPort = defaultPort
	}
	if config["shutdown_grace_time"] != "" {
		*flagShutdownGraceTime, err = strconv.Atoi(config["shutdown_grace_time"])
	} else {
		*flagShutdownGraceTime = defaultShutdownGraceTime
	}
	*flagConcurrency = config["concurrency"]
	return err
}

func parseConcurrency(value string) (map[string]int, error) {
	concurrency := map[string]int{}
	if strings.TrimSpace(value) == "" {
		return concurrency, nil
	}

	parts := strings.Split(value, ",")
	for _, part := range parts {
		if !strings.Contains(part, "=") {
			return concurrency, errors.New("Concurrency should be in the format: foo=1,bar=2")
		}

		nameValue := strings.Split(part, "=")
		n, v := strings.TrimSpace(nameValue[0]), strings.TrimSpace(nameValue[1])
		if n == "" || v == "" {
			return concurrency, errors.New("Concurrency should be in the format: foo=1,bar=2")
		}

		numProcs, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return concurrency, err
		}

		concurrency[n] = int(numProcs)
	}
	return concurrency, nil
}

type Forego struct {
	outletFactory *OutletFactory

	teardown, teardownNow Barrier // signal shutting down

	wg sync.WaitGroup
}

func (f *Forego) monitorInterrupt() {
	handler := make(chan os.Signal, 1)
	signal.Notify(handler, syscall.SIGALRM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	first := true

	for sig := range handler {
		switch sig {
		case syscall.SIGINT:
			fmt.Println("      | ctrl-c detected")
			fallthrough
		default:
			f.teardown.Fall()
			if !first {
				f.teardownNow.Fall()
			}
			first = false
		}
	}
}

func basePort(env Env) (int, error) {
	if flagPort != defaultPort {
		return flagPort, nil
	} else if env["PORT"] != "" {
		return strconv.Atoi(env["PORT"])
	} else if os.Getenv("PORT") != "" {
		return strconv.Atoi(os.Getenv("PORT"))
	}
	return defaultPort, nil
}

func (f *Forego) startProcess(idx, procNum int, proc ProcfileEntry, env Env, of *OutletFactory) {
	port, err := basePort(env)
	if err != nil {
		panic(err)
	}

	port = port + (idx * 100)

	const interactive = false
	workDir := filepath.Dir(flagProcfile)
	ps := NewProcess(workDir, proc.Command, env, interactive)
	procName := fmt.Sprint(proc.Name, ".", procNum+1)
	ps.Env["PORT"] = strconv.Itoa(port)

	ps.Stdin = nil

	stdout, err := ps.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderr, err := ps.StderrPipe()
	if err != nil {
		panic(err)
	}

	pipeWait := new(sync.WaitGroup)
	pipeWait.Add(2)
	go of.LineReader(pipeWait, procName, idx, stdout, false)
	go of.LineReader(pipeWait, procName, idx, stderr, true)

	of.SystemOutput(fmt.Sprintf("starting %s on port %d", procName, port))

	finished := make(chan struct{}) // closed on process exit

	err = ps.Start()
	if err != nil {
		f.teardown.Fall()
		of.SystemOutput(fmt.Sprint("Failed to start ", procName, ": ", err))
		return
	}

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		defer close(finished)
		pipeWait.Wait()
		ps.Wait()
	}()

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

		select {
		case <-finished:
			if flagRestart {
				f.startProcess(idx, procNum, proc, env, of)
			} else {
				f.teardown.Fall()
			}

		case <-f.teardown.Barrier():
			// Forego tearing down

			if !osHaveSigTerm {
				of.SystemOutput(fmt.Sprintf("Killing %s", procName))
				ps.Process.Kill()
				return
			}

			of.SystemOutput(fmt.Sprintf("sending SIGTERM to %s", procName))
			ps.SendSigTerm()

			// Give the process a chance to exit, otherwise kill it.
			select {
			case <-f.teardownNow.Barrier():
				of.SystemOutput(fmt.Sprintf("Killing %s", procName))
				ps.SendSigKill()
			case <-finished:
			}
		}
	}()
}

func runStart(cmd *Command, args []string) {
	pf, err := ReadProcfile(flagProcfile)
	handleError(err)

	concurrency, err := parseConcurrency(flagConcurrency)
	handleError(err)

	env, err := loadEnvs(envs)
	handleError(err)

	of := NewOutletFactory()
	of.Padding = pf.LongestProcessName(concurrency)

	f := &Forego{
		outletFactory: of,
	}

	go f.monitorInterrupt()

	// When teardown fires, start the grace timer
	f.teardown.FallHook = func() {
		go func() {
			time.Sleep(time.Duration(flagShutdownGraceTime) * time.Second)
			of.SystemOutput("Grace time expired")
			f.teardownNow.Fall()
		}()
	}

	var singleton string = ""
	if len(args) > 0 {
		singleton = args[0]
		if !pf.HasProcess(singleton) {
			of.ErrorOutput(fmt.Sprintf("no such process: %s", singleton))
		}
	}

	defaultConcurrency := 1

	var all bool
	for name, num := range concurrency {
		if name == "all" {
			defaultConcurrency = num
			all = true
		}
	}

	for idx, proc := range pf.Entries {
		numProcs := defaultConcurrency
		if len(concurrency) > 0 {
			if value, ok := concurrency[proc.Name]; ok {
				numProcs = value
			} else if !all {
				continue
			}
		}
		for i := 0; i < numProcs; i++ {
			if (singleton == "") || (singleton == proc.Name) {
				f.startProcess(idx, i, proc, env, of)
			}
		}
	}

	<-f.teardown.Barrier()

	f.wg.Wait()
}
