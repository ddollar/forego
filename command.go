package main

import (
  "flag"
  "fmt"
  "strings"
)

var flagEnv string
var flagProcfile string

type Command struct {
  // args does not include the command name
  Run  func(cmd *Command, args []string)
  Flag flag.FlagSet

  Usage string // first word is the command name
  Short string // `hk help` output
  Long  string // `hk help cmd` output
}

func (c *Command) printUsage() {
  if c.Runnable() {
    fmt.Printf("Usage: hk %s\n\n", c.Usage)
  }
  fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) Name() string {
  name := c.Usage
  i := strings.Index(name, " ")
  if i >= 0 {
    name = name[:i]
  }
  return name
}

func (c *Command) Runnable() bool {
  return c.Run != nil
}

func (c *Command) List() bool {
  return c.Short != ""
}
