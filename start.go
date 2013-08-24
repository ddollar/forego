package main

import (
  "fmt"
  "github.com/kr/pretty"
  "os"
  "os/exec"
  "path/filepath"
  "sync"
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

  ps_env := []string{}
  for name, val := range env {
    ps_env = append(ps_env, fmt.Sprintf("%s=%s", name, val))
  }
  ps_env = append(ps_env, "LANG=en_US.UTF-8")

  for idx, proc := range pf.entries {
    wg.Add(1)
    go func(idx int, proc ProcfileEntry) {
      command := []string{"/bin/bash", "-c", proc.command}
      ps := exec.Command(command[0], command[1:]...)
      ps.Dir = root
      ps.Env = ps_env
      ps.Stdin = nil
      ps.Stdout = createOutlet(proc.name, idx, false)
      ps.Stderr = createOutlet(proc.name, idx, true)
      ps.Start()
      ps.Wait()
      fmt.Println("process died")
      wg.Done()
    }(idx, proc)
  }
  wg.Wait()
  fmt.Println("afterdone")
}
