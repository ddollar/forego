package main

import (
	"os"
)

var commands = []*Command{
	cmdStart,
	cmdRun,
	cmdVersion,
	cmdHelp,
}

func RegisterCommand(c *Command) {
	commands = append(commands, c)
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.printUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			cmd.Run(cmd, cmd.Flag.Args())
			return
		}
	}
	usage()
}
