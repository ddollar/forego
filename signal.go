package main

import (
  "fmt"
  "os"
  "os/signal"
)

var _ = fmt.Println // lol

func init() {
  handler := make(chan os.Signal, 1)
  signal.Notify(handler, os.Interrupt)
  go func() {
    for sig := range handler {
      switch (sig) {
        case os.Interrupt:
          fmt.Println("      | ctrl-c detected")
          go func() { ShutdownProcesses() }()
        }
    }
  }()
}

