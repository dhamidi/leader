package main

import "fmt"

type ClearMenu struct {
	State *MenuState
}

func (cmd *ClearMenu) Execute() {
	for cmd.State.LinesDrawn > 0 {
		fmt.Fprintf(cmd.State.Out, "\033[2K\033[1A")
		cmd.State.LinesDrawn--
	}
}
