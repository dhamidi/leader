package main

import (
	"fmt"
	"io"
)

type MenuState struct {
	Out        io.Writer
	LinesDrawn int
	KeyMap     *KeyMap
}

type ClearMenu struct {
	State *MenuState
}

func (cmd *ClearMenu) Execute() {
	for cmd.State.LinesDrawn > 0 {
		fmt.Fprintf(cmd.State.Out, "\033[2K\033[1A")
		cmd.State.LinesDrawn--
	}
}

type DisplayMenu struct {
	State *MenuState
}

func (cmd *DisplayMenu) Execute() {
	clear := &ClearMenu{State: cmd.State}
	clear.Execute()
	fmt.Fprintf(cmd.State.Out, "\033[1B\r> %s\033[1B\r", cmd.State.KeyMap.Name)
	cmd.State.LinesDrawn += 2
	for key, value := range cmd.State.KeyMap.Keys {
		cmd.State.LinesDrawn++
		if child, isKeyMap := value.(*KeyMap); isKeyMap {
			fmt.Fprintf(cmd.State.Out, "[%c] %s\n\r", key, child.Name)
			continue
		}
		fmt.Fprintf(cmd.State.Out, "[%c] %s\n\r", key, value)
	}
}

type SelectMenuItem struct {
	State         *MenuState
	Key           rune
	BeforeExecute func()
	AfterExecute  func()
}

func (cmd *SelectMenuItem) Execute() {
	nextHandler, command := cmd.State.KeyMap.HandleKey(cmd.Key)
	if nextHandler == cmd.State.KeyMap && command == nil {
		return
	}
	if command != nil {
		cmd.BeforeExecute()
		command.Execute()
		cmd.AfterExecute()
	} else {
		cmd.State.KeyMap = nextHandler
	}
}
