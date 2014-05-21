package main

// +build darwin freebsd linux netbsd openbsd

import (
	"fmt"
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

func (p *Process) PlatformSpecificInit() {
	if !p.Interactive {
		p.cmd.SysProcAttr = &syscall.SysProcAttr{}
		p.cmd.SysProcAttr.Setsid = true
	}
	return
}

func (p *Process) SendSigTerm() {
	p.Signal(syscall.SIGTERM)
}

func (p *Process) SendSigKill() {
	p.Signal(syscall.SIGKILL)
}
