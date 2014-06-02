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

var cmdStart = &Command{
	Run:   runStart,
	Usage: "start [process name]... [-f procfile] [-e env] [-c concurrency] [-p port] [-r]",
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

type Forego struct {
	outletFactory *OutletFactory

	teardown, teardownNow Barrier // signal shutting down

	wg sync.WaitGroup
}

func (f *Forego) monitorInterrupt() {
	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)

	first := true

	for sig := range handler {
		switch sig {
		case os.Interrupt:
			fmt.Println("      | ctrl-c detected")

			f.teardown.Fall()
			if !first {
				f.teardownNow.Fall()
			}
			first = false
		}
	}
}

func (f *Forego) startProcess(idx, procNum int, proc ProcfileEntry, env Env, of *OutletFactory) {
	port := flagPort + (idx * 100)

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

	ps.Start()

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

		// Prevent goroutine from exiting before process has finished.
		defer func() { <-finished }()
		defer f.teardown.Fall()

		select {
		case <-finished:
			if flagRestart {
				f.startProcess(idx, procNum, proc, env, of)
				return
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

	f := &Forego{
		outletFactory: of,
	}

	go f.monitorInterrupt()

	// When teardown fires, start the grace timer
	f.teardown.FallHook = func() {
		go func() {
			time.Sleep(shutdownGraceTime)
			of.SystemOutput("Grace time expired")
			f.teardownNow.Fall()
		}()
	}

	var procsToRun = pf.Entries

	if len(args) > 0 {

		procsToRun = []ProcfileEntry{}
		for _, arg := range args {
			proc, ok := pf.GetProcess(arg)
			if !ok {
				of.ErrorOutput(fmt.Sprintf("Unknown proc '%s'", arg))
				return
			}
			procsToRun = append(procsToRun, proc)
		}
	}

	for idx, proc := range procsToRun {
		numProcs := 1
		if value, ok := concurrency[proc.Name]; ok {
			numProcs = value
		}
		for i := 0; i < numProcs; i++ {
			f.startProcess(idx, i, proc, env, of)
		}
	}

	<-f.teardown.Barrier()

	f.wg.Wait()
}
