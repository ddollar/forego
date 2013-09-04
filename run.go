package main

import (
	"os"
	"path/filepath"
	"strings"
)

var cmdRun = &Command{
	Run:   runRun,
	Usage: "run [-e env] [-c concurrency] [-p port]",
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
	if flagEnv == "" {
		root, _ := os.Getwd()
		flagEnv = filepath.Join(root, ".env")
	}

	env, err := ReadEnv(flagEnv)
	handleError(err)

	ps := NewProcess(strings.Join(args, " "), env)
	ps.Stdin = nil
	ps.Stdout = os.Stdout
	ps.Stderr = os.Stderr
	ps.Start()
	ps.Wait()
}
