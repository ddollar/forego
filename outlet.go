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
    formatter := fmt.Sprintf("%%-%ds | ", LongestOutletName())
    ct.ChangeColor(o.Color, true, ct.None, true)
    fmt.Printf(formatter, o.Name)
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

var outlets = map[string]*Outlet{}

func createOutlet(name string, index int, isError bool) *Outlet {
  outlets[name] = &Outlet{name, colors[index%len(colors)], isError}
  return outlets[name]
}

func LongestOutletName() (longest int) {
  // kr? better way to do this?
  longest = 0
  for name, _ := range outlets {
    if len(name) > longest {
      longest = len(name)
    }
  }
  return
}
