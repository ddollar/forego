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

	ps := NewProcess(strings.Join(args, " "), env)
	ps.Interactive = true
	ps.Root = workDir
	ps.Stdin = os.Stdin
	ps.Stdout = os.Stdout
	ps.Stderr = os.Stderr
	ps.Start()
	ps.Wait()
}
