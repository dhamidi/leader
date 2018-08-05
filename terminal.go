package main

import (
	"os"
	"syscall"

	"github.com/pkg/term/termios"
)

func rawTerminal() func() {
	state := new(syscall.Termios)
	if err := termios.Tcgetattr(os.Stdin.Fd(), state); err != nil {
		panic(err)
	}
	originalState := *state
	termios.Cfmakeraw(state)
	if err := termios.Tcsetattr(os.Stdin.Fd(), termios.TCSANOW, state); err != nil {
		panic(err)
	}
	return func() {
		termios.Tcsetattr(os.Stdin.Fd(), termios.TCSANOW, &originalState)
	}
}
