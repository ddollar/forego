package main

import (
  "fmt"
  "os"
  "os/signal"
)

var _ = fmt.Println // lol

func init() {
  handler := make(chan os.Signal, 1)
  fmt.Println("init")
  signal.Notify(handler, os.Interrupt)
  go func() {
    for sig := range handler {
      switch (sig) {
        case os.Interrupt:
          ShutdownProcesses()
        }
    }
  }()
}

