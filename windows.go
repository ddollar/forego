package main

// +build windows

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const osHaveSigTerm = false

func ShellInvocationCommand(interactive bool, root, command string) []string {
	return []string{"cmd", "/C", command}
}

func (p *Process) PlatformSpecificInit() {
	// NOP on windows for now.
	return
}

func (p *Process) Start() {
	command := ShellInvocationCommand(p.Root, p.Command)
	p.cmd = exec.Command(command[0], command[1:]...)
	p.cmd.Dir = p.Root
	p.cmd.Env = p.Env.asArray()
	p.cmd.Stdin = p.Stdin
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
	p.cmd.Start()
}

func (p *Process) SendSigTerm() {
	panic("SendSigTerm() not implemented on this platform")
}

func (p *Process) SendSigKill() {
	p.Signal(syscall.SIGKILL)
}
