package main

import (
  "fmt"
  "github.com/kr/pretty"
)

var flagProcfile string
var flagEnv string

var cmdStart = &Command{
  Run: runStart,
  Usage: "start [-f procfile] [-e env] [-c concurrency]",
  Short: "start the app",
  Long: `
Start the application specified by a Procfile (defaults to ./Procfile)

Examples:

    forego start

    forego start -f Procfile.test -e .env.test
`,
}

func init() {
  cmdStart.Flag.StringVar(&flagProcfile, "f", "./Procfile", "procfile")
  cmdStart.Flag.StringVar(&flagEnv, "e", "", "env")
}

func runStart(cmd *Command, args []string) {
  pf, err := OpenProcfile(flagProcfile)
  if err != nil {
    fmt.Println("ERROR:", err)
    return
  }
  pretty.Println("procfile", flagProcfile, pf, err)
}
