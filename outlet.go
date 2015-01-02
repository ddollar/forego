package main

import (
	"bufio"
	"fmt"
	"github.com/ddollar/forego/Godeps/_workspace/src/github.com/daviddengcn/go-colortext"
	"io"
	"os"
	"sync"
	"bytes"
)

type OutletFactory struct {
	Padding int

	sync.Mutex
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

func (of *OutletFactory) LineReader(wg *sync.WaitGroup, name string, index int, r io.Reader, isError bool) {
	defer wg.Done()

	color := colors[index%len(colors)]

	reader := bufio.NewReader(r)

	var buffer bytes.Buffer

	for {
		buf := make([]byte, 1024)
		v, _ := reader.Read(buf)

		if v == 0 {
			return
		}

		idx := bytes.IndexByte(buf, '\n')
		if idx >= 0 {
			buffer.Write(buf[0:idx])
			of.WriteLine(name, buffer.String(), color, ct.None, isError)
			buffer.Reset()
		} else {
			buffer.Write(buf)
		}
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
