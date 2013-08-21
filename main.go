package main

import (
  "github.com/kr/pretty"
)

func main() {
  pf, err := OpenProcfile("example/Procfile")

  pretty.Println("testing", pf, err);
}
