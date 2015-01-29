package main

import (
	"os"
	"strings"
)

var cmdRun = &Command{
	Run:   runRun,
	Usage: "run [-e env] [-p port]",
	Short: "Run a one-off command",
	Long: `
Run a one-off command

Examples:

  forego run bin/migrate
`,
}

var runEnvs envFiles

func init() {
	cmdRun.Flag.Var(&runEnvs, "e", "env")
}

func runRun(cmd *Command, args []string) {
	if len(args) < 1 {
		cmd.printUsage()
		os.Exit(1)
	}
	workDir, err := os.Getwd()
	if err != nil {
		handleError(err)
	}

	env, err := loadEnvs(runEnvs, "", "")
	handleError(err)

	const interactive = true
	ps := NewProcess(workDir, strings.Join(args, " "), env, interactive)
	ps.Stdin = os.Stdin
	ps.Stdout = os.Stdout
	ps.Stderr = os.Stderr

	err = ps.Start()
	handleError(err)

	err = ps.Wait()
	handleError(err)
}
