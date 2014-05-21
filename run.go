package main

import (
	"os"
	"path/filepath"
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

func init() {
	cmdRun.Flag.StringVar(&flagEnv, "e", ".env", "env")
}

func runRun(cmd *Command, args []string) {
	workDir, err := os.Getwd()
	if err != nil {
		handleError(err)
	}
	if flagEnv == "" {
		flagEnv = filepath.Join(workDir, ".env")
	}

	env, err := ReadEnv(flagEnv)
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
