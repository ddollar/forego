package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io"
	"os"
	"sync"
)

type OutletFactory struct {
	Outlets map[string]*Outlet
	Padding int
}

type Outlet struct {
	Name    string
	Color   ct.Color
	IsError bool
	Factory *OutletFactory
}

var mx sync.Mutex

var colors = []ct.Color{
	ct.Cyan,
	ct.Yellow,
	ct.Green,
	ct.Magenta,
	ct.Red,
	ct.Blue,
}

func NewOutletFactory() (of *OutletFactory) {
	of = new(OutletFactory)
	of.Outlets = make(map[string]*Outlet)
	return
}

func (o *Outlet) Write(b []byte) (num int, err error) {
	mx.Lock()
	defer mx.Unlock()
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		formatter := fmt.Sprintf("%%-%ds | ", o.Factory.Padding)
		ct.ChangeColor(o.Color, true, ct.None, false)
		fmt.Printf(formatter, o.Name)
		if o.IsError {
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

func (of *OutletFactory) CreateOutlet(name string, index int, isError bool) *Outlet {
	of.Outlets[name] = &Outlet{name, colors[index%len(colors)], isError, of}
	return of.Outlets[name]
}

func (of *OutletFactory) SystemOutput(str string) {
	ct.ChangeColor(ct.White, true, ct.None, false)
	formatter := fmt.Sprintf("%%-%ds | ", of.Padding)
	fmt.Printf(formatter, "forego")
	ct.ResetColor()
	fmt.Println(str)
	ct.ResetColor()
}

func (of *OutletFactory) ErrorOutput(str string) {
	fmt.Printf("ERROR: %s\n", str)
	os.Exit(1)
}
