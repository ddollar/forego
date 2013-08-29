package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Process struct {
	Command string
	Env     Env
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Root    string

	cmd *exec.Cmd
}

func NewProcess(command string, env Env) (p *Process) {
	p = new(Process)
	p.Command = command
	p.Env = env
	p.Stdin = os.Stdin
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Root, _ = os.Getwd()
	return
}

func (p *Process) Start() {
	command := []string{"/bin/bash", "-c", fmt.Sprintf("source \"%s\" 2>/dev/null; %s", filepath.Join(p.Root, ".profile"), p.Command)}
	p.cmd = exec.Command(command[0], command[1:]...)
	p.cmd.Dir = p.Root
	p.cmd.Env = p.envAsArray()
	p.cmd.Stdin = p.Stdin
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
	p.cmd.SysProcAttr = &syscall.SysProcAttr{}
	p.cmd.SysProcAttr.Setsid = true
	p.cmd.Start()
}

func (p *Process) Wait() {
  p.cmd.Wait()
}

func (p *Process) Signal(signal syscall.Signal) {
	group, _ := os.FindProcess(-1 * p.cmd.Process.Pid)
	group.Signal(signal)
}

func (p *Process) envAsArray() (env []string) {
	for name, val := range os.Environ() {
		env = append(env, fmt.Sprintf("%s=%s", name, val));
	}
	for name, val := range p.Env {
		env = append(env, fmt.Sprintf("%s=%s", name, val));
	}
	return
}

