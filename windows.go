//go:build windows
// +build windows

package main

import (
	"syscall"
)

const osHaveSigTerm = false

func ShellInvocationCommand(interactive bool, root, command string) []string {
	return []string{"cmd", "/C", command}
}

func (p *Process) PlatformSpecificInit() {
	p.Dir = p.WorkDir
	// cmd.exe expects its command string verbatim, not escaped as a Windows argv.
	p.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: syscall.EscapeArg(p.Path) + " /S /C \"" + p.Command + "\"",
	}
	return
}

func (p *Process) SendSigTerm() {
	panic("SendSigTerm() not implemented on this platform")
}

func (p *Process) SendSigKill() {
	p.Signal(syscall.SIGKILL)
}
