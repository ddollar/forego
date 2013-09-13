package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Process struct {
	Command     string
	Env         Env
	Interactive bool
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	Root        string

	cmd *exec.Cmd
}

func NewProcess(command string, env Env) (p *Process) {
	p = new(Process)
	p.Command = command
	p.Env = env
	p.Interactive = false
	p.Stdin = os.Stdin
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Root, _ = os.Getwd()
	return
}

func (p *Process) Wait() {
	p.cmd.Wait()
}

func (p *Process) shellArgument() string {
	if p.Interactive {
		return "-ic"
	} else {
		return "-c"
	}
}

func (p *Process) envAsArray() (env []string) {
	for _, pair := range os.Environ() {
		env = append(env, pair)
	}
	for name, val := range p.Env {
		env = append(env, fmt.Sprintf("%s=%s", name, val))
	}
	return
}
