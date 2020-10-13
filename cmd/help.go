package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

var cmdHelp = &Command{
	Usage: "help [topic]",
	Short: "Show this help",
	Long:  `Help shows usage for a command.`,
}

func init() {
	cmdHelp.Run = runHelp // break init loop
}

func runHelp(cmd *Command, args []string) {
	if len(args) == 0 {
		printUsage()
		return
	}
	if len(args) != 1 {
		log.Fatal("too many arguments")
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.printUsage()
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic: %q. Run 'forego help'.\n", args[0])
	os.Exit(2)
}

var usageTemplate = template.Must(template.New("usage").Parse(`
Usage: forego <command> [<args>]

Available commands:{{range .Commands}}{{if .Runnable}}{{if .List}}
   {{.Name | printf "%-8s"}}  {{.Short}}{{end}}{{end}}{{end}}

Run 'forego help [command]' for details.
`[1:]))

func printUsage() {
	usageTemplate.Execute(os.Stdout, struct {
		Commands []*Command
	}{
		commands,
	})
}

func usage() {
	printUsage()
	os.Exit(2)
}
