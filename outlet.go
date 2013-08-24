package main

import (
  "bufio"
  "bytes"
  "fmt"
  "github.com/daviddengcn/go-colortext"
  "github.com/kr/pretty"
  "sync"
)

type Outlet struct {
  Name string
  Color ct.Color
  IsError bool
}

var _ = pretty.Println // lol
var _ = bufio.NewScanner
var _ = bytes.NewReader

var mutex = new(sync.Mutex)

var colors = []ct.Color{
  ct.Cyan,
  ct.Yellow,
  ct.Green,
  ct.Magenta,
  ct.Red,
  ct.Blue,
}

func (o *Outlet) Write(b []byte) (num int, err error) {
  mutex.Lock()
  defer mutex.Unlock()
  scanner := bufio.NewScanner(bytes.NewReader(b))
  for scanner.Scan() {
    ct.ChangeColor(o.Color, true, ct.None, true)
    fmt.Printf("%-10s | ", o.Name)
    if (o.IsError) {
      ct.ChangeColor(ct.Red, true, ct.None, true)
    } else {
      ct.ResetColor()
    }
    fmt.Println(scanner.Text())
    ct.ResetColor()
  }
  num = len(b)
  return
}

func createOutlet(name string, index int, isError bool) *Outlet {
  return &Outlet{name, colors[index%len(colors)], isError}
}
