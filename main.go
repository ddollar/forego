package main

import "os"

var commands = []*Command{
	cmdStart,
	cmdRun,
	// cmdUpdate,
	cmdVersion,
	cmdHelp,
}

var allowUpdate string = "true"

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
	}

	if allowUpdate == "false" {
		cmdUpdate.Disabled = true
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
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
