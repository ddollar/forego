package main

import (
	"fmt"
	"github.com/kr/pretty"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var _ = pretty.Println // lol
var _ = os.Stdout

var cmdRun = &Command{
	Run:   runRun,
	Usage: "run [-e env] [-c concurrency] [-p port]",
	Short: "Run a one-off command",
	Long: `
Run a one-off command

Examples:

  forego run bin/migrate
`,
}

func init() {
	cmdRun.Flag.StringVar(&flagEnv, "e", ".env", "env")
}

func runRun(cmd *Command, args []string) {
	command := []string{"/bin/bash", "-c"}
	command = append(command, fmt.Sprintf("source .profile 2>/dev/null; %s", strings.Join(args, " ")))

	env, err := ReadEnv(flagEnv)
	handleError(err)

	ps_env := os.Environ()
	for name, val := range env {
		ps_env = append(ps_env, fmt.Sprintf("%s=%s", name, val))
	}

	ps := exec.Command(command[0], command[1:]...)
	ps.Dir, _ = os.Getwd()
	ps.Env = ps_env
	ps.Stdin = nil
	ps.Stdout = os.Stdout
	ps.Stderr = os.Stderr
	ps.SysProcAttr = &syscall.SysProcAttr{}
	ps.SysProcAttr.Setsid = true
	ps.Start()
	ps.Wait()
}
