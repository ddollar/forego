package main

// +build darwin freebsd linux netbsd openbsd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const osHaveSigTerm = true

func (p *Process) Start() {
	command := []string{"/bin/bash", p.shellArgument(), fmt.Sprintf("source \"%s\" 2>/dev/null; %s", filepath.Join(p.Root, ".profile"), p.Command)}
	p.cmd = exec.Command(command[0], command[1:]...)
	p.cmd.Dir = p.Root
	p.cmd.Env = p.envAsArray()
	p.cmd.Stdin = p.Stdin
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
	if !p.Interactive {
		p.cmd.SysProcAttr = &syscall.SysProcAttr{}
		p.cmd.SysProcAttr.Setsid = true
	}
	p.cmd.Start()
}

func (p *Process) Signal(signal syscall.Signal) {
	if p.Running() {
		group, _ := os.FindProcess(-1 * p.Pid())
		group.Signal(signal)
	}
}

func (p *Process) SendSigTerm() {
	p.Signal(syscall.SIGTERM)
}

func (p *Process) SendSigKill() {
	p.Signal(syscall.SIGKILL)
}
