package catcher

import (
	"os"
	"os/signal"
)

// Signal is in charge of catching signals for the application.
type Signal struct {
	c chan os.Signal
}

// NewSignal will create a new signal cather for the given signals.
func NewSignal(signals ...os.Signal) *Signal {
	sig := &Signal{
		c: make(chan os.Signal),
	}
	signal.Notify(sig.c, signals...)
	return sig
}

// Wait will wait for one of the watched signals to happen.
func (sig *Signal) Wait() os.Signal {
	s := <-sig.c
	signal.Stop(sig.c)
	close(sig.c)
	return s
}
