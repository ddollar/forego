// +build darwin freebsd linux netbsd openbsd

package main

import (
	"fmt"
	"syscall"
)

const osHaveSigTerm = true

func ShellInvocationCommand(interactive bool, root, command string) []string {
	shellArgument := "-c"
	if interactive {
		shellArgument = "-ic"
	}
	shellCommand := fmt.Sprintf("cd \"%s\"; source .profile 2>/dev/null; exec %s", root, command)
	return []string{"bash", shellArgument, shellCommand}
}

func (p *Process) PlatformSpecificInit() {
	if !p.Interactive {
		p.SysProcAttr = &syscall.SysProcAttr{}
		p.SysProcAttr.Setsid = true
	}
	return
}

func (p *Process) SendSigTerm() {
	p.Signal(syscall.SIGTERM)
}

func (p *Process) SendSigKill() {
	p.Signal(syscall.SIGKILL)
}
