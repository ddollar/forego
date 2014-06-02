package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	print("Foo")
	time.Sleep(10 * time.Millisecond)
	println("Bar")

	print("Baz")
	time.Sleep(10 * time.Millisecond)
	println("Qux")

	fmt.Fprintln(os.Stdout, "This is on \x1b[32mstdout")

	os.Stdout.Close()

	s := rand.Intn(3) + 1

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		if rand.Intn(4) == 1 {
			println("IGNORING EXIT")
			time.Sleep(100 * time.Second)
		}
		println("Got SIGTERM")
	case <-time.After(time.Duration(s) * time.Second):
		println("Timed out")
	}

}
