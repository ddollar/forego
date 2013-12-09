package main

// +build windows

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func (p *Process) Start() {
	command := []string{"cmd", "/C", p.Command}
	p.cmd = exec.Command(command[0], command[1:]...)
	p.cmd.Dir = p.Root
	p.cmd.Env = p.envAsArray()
	p.cmd.Stdin = p.Stdin
	p.cmd.Stdout = p.Stdout
	p.cmd.Stderr = p.Stderr
	p.cmd.Start()
}

func (p *Process) Signal(signal syscall.Signal) {
	group, _ := os.FindProcess(-1 * p.cmd.Process.Pid)
	group.Signal(signal)
}

func ShutdownProcesses(of *OutletFactory) {
	shutdown_mutex.Lock()
	of.SystemOutput("shutting down")
	for name, ps := range processes {
		of.SystemOutput(fmt.Sprintf("terminating %s", name))
		ps.cmd.Process.Signal(os.Kill)
	}
	os.Exit(1)
}
