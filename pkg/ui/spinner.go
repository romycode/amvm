package ui

import (
	"fmt"
	"time"
)

var spinChars = `|/-\`

type Spinner struct {
	stop chan bool
	done chan bool
	msg  string
	i    int
}

func NewSpinner(msg string) *Spinner {
	return &Spinner{make(chan bool), make(chan bool), msg, 0}
}

func (r *Spinner) Start() {
	go func() {
		for {
			select {
			case <-r.stop:
				fmt.Printf("\r\033[K")
				break
			case <-time.After(250 * time.Millisecond):
				fmt.Printf("%s %c \r", r.msg, spinChars[r.i])
				r.i = (r.i + 1) % len(spinChars)
			}
		}
	}()
}

func (r *Spinner) Stop() {
	r.stop <- true
}
