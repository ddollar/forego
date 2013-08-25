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
  "unsafe"
)

var _ = pretty.Println // lol
var _ = os.Stdout

var flagProcfile string
var flagEnv string
var wg sync.WaitGroup

var cmdStart = &Command{
  Run: runStart,
  Usage: "start [-f procfile] [-e env] [-c concurrency]",
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
    go func(idx int, proc ProcfileEntry) {
      command := []string{"/bin/bash", "-c", proc.command}
      var attr syscall.ProcAttr
      attr.Files = []uintptr{0,0,0}
      attr.Files[1] = uintptr(unsafe.Pointer(createOutlet(proc.name, idx, false)))
      attr.Files[2] = uintptr(unsafe.Pointer(createOutlet(proc.name, idx, true)))
      attr.Env = ps_env
      pid, err := syscall.ForkExec(command[0], command[1:], &attr)
      pretty.Println("arg0", command[0], "argv", command[1:])
      fmt.Println("pid", pid, "err", err)
      /* ps := exec.Command(command[0], command[1:]...)*/
      /* ps.Dir = root*/
      /* ps.Env = ps_env*/
      /* ps.Stdin = nil*/
      /* ps.Stdout = createOutlet(proc.name, idx, false)*/
      /* ps.Stderr = createOutlet(proc.name, idx, true)*/
      /* processes[proc.name] = ps*/
      /* ps.Start()*/
      /* syscall.Setpgid(ps.Process.Pid, -1 * ps.Process.Pid)*/
      fmt.Println("pid", pid)
      var wait syscall.WaitStatus
      var usage syscall.Rusage
      syscall.Wait4(pid, &wait, 0, &usage)
      fmt.Println("exitstatus", wait.ExitStatus())
      fmt.Println("process died")
      wg.Done()
      delete(processes, proc.name)
    }(idx, proc)
  }
  wg.Wait()
  fmt.Println("afterdone")
}

func ShutdownProcesses() {
  for _, ps := range processes {
    fmt.Println("killing", ps.Process)
    group, _ := os.FindProcess(ps.Process.Pid)
    group.Signal(syscall.SIGTERM)
  }
  go func() {
    time.Sleep(1 * time.Second)
    fmt.Println("really killing")
    for _, ps := range processes {
      fmt.Println("killing", ps.Process)
      ps.Process.Signal(syscall.SIGKILL)
      b, _ := os.FindProcess(-1 * ps.Process.Pid)
      b.Signal(syscall.SIGKILL)
    }
  }()
}
