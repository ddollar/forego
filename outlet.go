package main

import (
  "bufio"
  "bytes"
  "fmt"
  "github.com/daviddengcn/go-colortext"
  "github.com/kr/pretty"
  "io"
  "os"
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
    ct.ChangeColor(o.Color, true, ct.None, false)
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

func ProcessOutput(w io.Writer, str string) {
  w.Write([]byte(str))
}

var outlets = map[string]*Outlet{}

func createOutlet(name string, index int, isError bool) *Outlet {
  outlets[name] = &Outlet{name, colors[index%len(colors)], isError}
  return outlets[name]
}

func LongestOutletName() (longest int) {
  // kr?
  longest = 6 // forego is the shortest name
  for name, _ := range outlets {
    if len(name) > longest {
      longest = len(name)
    }
  }
  return
}

func SystemOutput(str string) {
  ct.ChangeColor(ct.White, true, ct.None, false)
  formatter := fmt.Sprintf("%%-%ds | ", LongestOutletName())
  fmt.Printf(formatter, "forego")
  ct.ResetColor()
  fmt.Println(str)
  ct.ResetColor()
}

func ErrorOutput(str string) {
  fmt.Printf("ERROR: %s\n", str)
  os.Exit(1)
}
