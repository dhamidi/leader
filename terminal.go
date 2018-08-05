package main

import (
	"os"

	"github.com/Nerdmaster/terminal"
)

func rawTerminal() func() {
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	return func() {
		terminal.Restore(int(os.Stdin.Fd()), oldState)
	}
}
