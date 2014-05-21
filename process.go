package main

import (
	"os"
	"os/exec"
	"syscall"
)

type Process struct {
	Command     string
	Env         Env
	Interactive bool

	*exec.Cmd
}

func NewProcess(workdir, command string, env Env, interactive bool) (p *Process) {
	argv := ShellInvocationCommand(interactive, workdir, command)
	return &Process{
		command, env, interactive, exec.Command(argv[0], argv[1:]...),
	}
}

func (p *Process) Start() error {
	p.Cmd.Env = p.Env.asArray()
	p.PlatformSpecificInit()
	return p.Cmd.Start()
}

func (p *Process) Signal(signal syscall.Signal) {
	group, _ := os.FindProcess(-1 * p.Process.Pid)
	group.Signal(signal)
}
