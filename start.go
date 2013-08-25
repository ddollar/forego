package main

import (
  "fmt"
  "github.com/kr/pretty"
  "os"
  "os/exec"
  "path/filepath"
  "sync"
  "syscall"
  "time"
)

const shutdownGraceTime = 1 * time.Second

var _ = pretty.Println // lol
var _ = os.Stdout

var flagProcfile string
var flagEnv string
var flagPort int

var wg sync.WaitGroup

var cmdStart = &Command{
  Run: runStart,
  Usage: "start [-f procfile] [-e env] [-c concurrency] [-p port]",
  Short: "Start the application",
  Long: `
Start the application specified by a Procfile (defaults to ./Procfile)

Examples:

  forego start
  forego start -f Procfile.test -e .env.test
`,
}

var processes = map[string]*exec.Cmd{}

func init() {
  cmdStart.Flag.StringVar(&flagProcfile, "f", "Procfile", "procfile")
  cmdStart.Flag.StringVar(&flagEnv, "e", "", "env")
  cmdStart.Flag.IntVar(&flagPort, "p", 5000, "port")
}

func runStart(cmd *Command, args []string) {
  root := filepath.Dir(flagProcfile)
  if (flagEnv == "") {
    flagEnv = filepath.Join(root, ".env")
  }
  pf, err := ReadProcfile(flagProcfile)
  handleError(err)
  env, err := ReadEnv(flagEnv)
  handleError(err)

  ps_env := os.Environ()
  for name, val := range env {
    ps_env = append(ps_env, fmt.Sprintf("%s=%s", name, val))
  }

  for idx, proc := range pf.entries {
    wg.Add(1)
    command := []string{"/bin/bash", "-c", proc.Command}
    ps := exec.Command(command[0], command[1:]...)
    port := flagPort + (idx * 100)
    processes[proc.Name] = ps
    ps.Dir = root
    ps.Env = append(ps_env, fmt.Sprintf("PORT=%d", port))
    ps.Stdin = nil
    ps.Stdout = createOutlet(proc.Name, idx, false)
    ps.Stderr = createOutlet(proc.Name, idx, true)
    ps.SysProcAttr = &syscall.SysProcAttr{}
    ps.SysProcAttr.Setsid = true
    ps.Start()
    SystemOutput(fmt.Sprintf("starting %s on port %d", proc.Name, port))
    go func(proc ProcfileEntry, ps *exec.Cmd) {
      ps.Wait()
      wg.Done()
      delete(processes, proc.Name)
    }(proc, ps)
  }
  wg.Wait()
}

func ShutdownProcesses() {
  for name, ps := range processes {
    SystemOutput(fmt.Sprintf("sending SIGTERM to %s", name))
    group, _ := os.FindProcess(-1 * ps.Process.Pid)
    group.Signal(syscall.SIGTERM)
  }
  go func() {
    time.Sleep(shutdownGraceTime)
    for name, ps := range processes {
      SystemOutput(fmt.Sprintf("sending SIGKILL to %s", name))
      group, _ := os.FindProcess(-1 * ps.Process.Pid)
      group.Signal(syscall.SIGKILL)
    }
  }()
}
