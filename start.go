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
	"time"
)

const shutdownGraceTime = 3 * time.Second

var flagPort int
var flagConcurrency string
var flagRestart bool
var shutdownNow bool

var processes = map[string]*Process{}
var shutdown_mutex = new(sync.Mutex)
var wg sync.WaitGroup

var cmdStart = &Command{
	Run:   runStart,
	Usage: "start [process name] [-f procfile] [-e env] [-c concurrency] [-p port] [-r]",
	Short: "Start the application",
	Long: `
Start the application specified by a Procfile (defaults to ./Procfile)

Examples:

  forego start
  forego start web
  forego start -f Procfile.test -e .env.test
`,
}

func init() {
	cmdStart.Flag.StringVar(&flagProcfile, "f", "Procfile", "procfile")
	cmdStart.Flag.StringVar(&flagEnv, "e", "", "env")
	cmdStart.Flag.IntVar(&flagPort, "p", 5000, "port")
	cmdStart.Flag.StringVar(&flagConcurrency, "c", "", "concurrency")
	cmdStart.Flag.BoolVar(&flagRestart, "r", false, "restart")
}

func parseConcurrency(value string) (map[string]int, error) {
	concurrency := map[string]int{}
	if strings.TrimSpace(value) == "" {
		return concurrency, nil
	}

	parts := strings.Split(value, ",")
	for _, part := range parts {
		if !strings.Contains(part, "=") {
			return concurrency, errors.New("Parsing concurency")
		}

		nameValue := strings.Split(part, "=")
		n, v := strings.TrimSpace(nameValue[0]), strings.TrimSpace(nameValue[1])
		if n == "" || v == "" {
			return concurrency, errors.New("Parsing concurency")
		}

		numProcs, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return concurrency, err
		}

		concurrency[n] = int(numProcs)
	}
	return concurrency, nil
}

func startProcess(idx, procNum int, proc ProcfileEntry, env Env, of *OutletFactory) {
	shutdown_mutex.Lock()
	wg.Add(1)
	port := flagPort + (idx * 100)
	ps := NewProcess(proc.Command, env)
	procName := strings.Join([]string{
		proc.Name,
		strconv.FormatInt(int64(procNum+1), 10)}, ".")
	processes[procName] = ps
	ps.Env["PORT"] = strconv.Itoa(port)
	ps.Root = filepath.Dir(flagProcfile)
	ps.Stdin = nil
	ps.Stdout = of.CreateOutlet(procName, idx, false)
	ps.Stderr = of.CreateOutlet(procName, idx, true)
	ps.Start()
	of.SystemOutput(fmt.Sprintf("starting %s on port %d", procName, port))
	go func(proc ProcfileEntry, ps *Process) {
		ps.Wait()

		if flagRestart && !shutdownNow {
			delete(processes, proc.Name)
			startProcess(idx, procNum, proc, env, of)
			wg.Done()
			return
		}

		wg.Done()
		delete(processes, proc.Name)
		ShutdownProcesses(of)

	}(proc, ps)
	shutdown_mutex.Unlock()
}

func runStart(cmd *Command, args []string) {
	root := filepath.Dir(flagProcfile)

	if flagEnv == "" {
		flagEnv = filepath.Join(root, ".env")
	}

	pf, err := ReadProcfile(flagProcfile)
	handleError(err)

	env, err := ReadEnv(flagEnv)
	handleError(err)

	concurrency, err := parseConcurrency(flagConcurrency)
	handleError(err)

	of := NewOutletFactory()
	of.Padding = pf.LongestProcessName(concurrency)

	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)

	go func() {
		for sig := range handler {
			switch sig {
			case os.Interrupt:
				shutdownNow = true
				fmt.Println("      | ctrl-c detected")
				go func() { ShutdownProcesses(of) }()
			}
		}
	}()

	var singleton string = ""
	if len(args) > 0 {
		singleton = args[0]
		if !pf.HasProcess(singleton) {
			of.ErrorOutput(fmt.Sprintf("no such process: %s", singleton))
		}
	}

	for idx, proc := range pf.Entries {
		numProcs := 1
		if value, ok := concurrency[proc.Name]; ok {
			numProcs = value
		}
		for i := 0; i < numProcs; i++ {
			if (singleton == "") || (singleton == proc.Name) {
				startProcess(idx, i, proc, env, of)
			}
		}
	}

	wg.Wait()
}
