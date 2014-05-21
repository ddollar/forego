package main

// +build darwin freebsd linux netbsd openbsd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
)

const osHaveSigTerm = true

func ShellInvocationCommand(interactive bool, root, command string) []string {
	shellArgument := "-c"
	if interactive {
		shellArgument = "-ic"
	}
	profile := filepath.Join(root, ".profile")
	shellCommand := fmt.Sprintf("source \"%s\" 2>/dev/null; %s", profile, command)
	return []string{"/bin/bash", shellArgument, shellCommand}

}

func (p *Process) Start() {
	command := ShellInvocationCommand(p.Interactive, p.Root, p.Command)
	p.cmd = exec.Command(command[0], command[1:]...)
	p.cmd.Dir = p.Root
	p.cmd.Env = p.Env.asArray()
	p.cmd.Stdin = p.Stdin
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
	if !p.Interactive {
		p.cmd.SysProcAttr = &syscall.SysProcAttr{}
		p.cmd.SysProcAttr.Setsid = true
	}
	p.cmd.Start()
}

func (p *Process) SendSigTerm() {
	p.Signal(syscall.SIGTERM)
}

func (p *Process) SendSigKill() {
	p.Signal(syscall.SIGKILL)
}
