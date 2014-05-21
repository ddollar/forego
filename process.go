package main

import (
	"io"
	"os"
	"os/exec"
	"syscall"
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

func (p *Process) Signal(signal syscall.Signal) {
	group, _ := os.FindProcess(-1 * p.Pid())
	group.Signal(signal)
}

func (p *Process) Pid() int {
	return p.cmd.Process.Pid
}

func (p *Process) Wait() {
	p.cmd.Wait()
}
