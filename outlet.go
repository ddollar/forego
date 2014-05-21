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
	Padding int

	sync.Mutex
}

type Outlet struct {
	Name    string
	Color   ct.Color
	IsError bool
	Factory *OutletFactory
}

var colors = []ct.Color{
	ct.Cyan,
	ct.Yellow,
	ct.Green,
	ct.Magenta,
	ct.Red,
	ct.Blue,
}

func NewOutletFactory() (of *OutletFactory) {
	return new(OutletFactory)
}

func (o *Outlet) Write(b []byte) (num int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		o.Factory.WriteLine(o.Name, scanner.Text(), ct.White, ct.None, o.IsError)
	}
	num = len(b)
	return
}

func ProcessOutput(w io.Writer, str string) {
	w.Write([]byte(str))
}

func (of *OutletFactory) LineReader(wg *sync.WaitGroup, name string, index int, r io.Reader, isError bool) {
	defer wg.Done()

	o := &Outlet{name, colors[index%len(colors)], isError, of}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		of.WriteLine(o.Name, scanner.Text(), o.Color, ct.None, o.IsError)
	}
}

func (of *OutletFactory) SystemOutput(str string) {
	of.WriteLine("forego", str, ct.White, ct.None, false)
}

func (of *OutletFactory) ErrorOutput(str string) {
	fmt.Printf("ERROR: %s\n", str)
	os.Exit(1)
}

// Write out a single coloured line
func (of *OutletFactory) WriteLine(left, right string, leftC, rightC ct.Color, isError bool) {
	of.Lock()
	defer of.Unlock()

	ct.ChangeColor(leftC, true, ct.None, false)
	formatter := fmt.Sprintf("%%-%ds | ", of.Padding)
	fmt.Printf(formatter, left)

	if isError {
		ct.ChangeColor(ct.Red, true, ct.None, true)
	} else {
		ct.ResetColor()
	}
	fmt.Println(right)
	if isError {
		ct.ResetColor()
	}
}
