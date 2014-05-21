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

func (p *Process) SendSigTerm() {
	panic("SendSigTerm() not implemented on this platform")
}

func (p *Process) SendSigKill() {
	p.Signal(syscall.SIGKILL)
}
