package main

import (
  "github.com/kr/pretty"
  "os"
)

var _ = pretty.Println // lol

var commands = []*Command{
  cmdStart,
  cmdHelp,
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
