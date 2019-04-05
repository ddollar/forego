package main

import (
	"sync"
)

// Direct import of https://github.com/pwaller/barrier/blob/master/barrier.go

// The zero of Barrier is a ready-to-use value
type Barrier struct {
	channel          chan struct{}
	fall, initialize sync.Once
	FallHook         func()
}

func (b *Barrier) init() {
	b.initialize.Do(func() { b.channel = make(chan struct{}) })
}

// `b.Fall()` can be called any number of times and causes the channel returned
// by `b.Barrier()` to become closed (permanently available for immediate reading)
func (b *Barrier) Fall() {
	b.init()
	b.fall.Do(func() {
		if b.FallHook != nil {
			b.FallHook()
		}
		close(b.channel)
	})
}

// When `b.Fall()` is called, the channel returned by Barrier() is closed
// (and becomes always readable)
func (b *Barrier) Barrier() <-chan struct{} {
	b.init()
	return b.channel
}
