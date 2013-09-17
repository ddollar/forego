package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const shutdownGraceTime = 3 * time.Second

var flagPort int

var processes = map[string]*Process{}
var shutdown_mutex = new(sync.Mutex)
var wg sync.WaitGroup

var cmdStart = &Command{
	Run:   runStart,
	Usage: "start [process name] [-f procfile] [-e env] [-c concurrency] [-p port]",
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

	of := NewOutletFactory()
	of.Padding = pf.LongestProcessName()

	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)

	go func() {
		for sig := range handler {
			switch sig {
			case os.Interrupt:
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
		if (singleton == "") || (singleton == proc.Name) {
			shutdown_mutex.Lock()
			wg.Add(1)
			port := flagPort + (idx * 100)
			ps := NewProcess(proc.Command, env)
			processes[proc.Name] = ps
			ps.Env["PORT"] = strconv.Itoa(flagPort + (idx * 1000))
			ps.Root = filepath.Dir(flagProcfile)
			ps.Stdin = nil
			ps.Stdout = of.CreateOutlet(proc.Name, idx, false)
			ps.Stderr = of.CreateOutlet(proc.Name, idx, true)
			ps.Start()
			of.SystemOutput(fmt.Sprintf("starting %s on port %d", proc.Name, port))
			go func(proc ProcfileEntry, ps *Process) {
				ps.Wait()
				wg.Done()
				delete(processes, proc.Name)
				ShutdownProcesses(of)
			}(proc, ps)
			shutdown_mutex.Unlock()
		}
	}

	wg.Wait()
}
