package main

import (
	"fmt"
	"sort"
)

type DisplayMenu struct {
	State *MenuState
}

func (cmd *DisplayMenu) Execute() {
	clear := &ClearMenu{State: cmd.State}
	clear.Execute()
	fmt.Fprintf(cmd.State.Out, "\033[2K> %s\n\r", cmd.State.CurrentPath())
	cmd.State.LinesDrawn++
	keys := []string{}
	for key := range cmd.State.CurrentHandler().Keys {
		keys = append(keys, fmt.Sprintf("%c", key))
	}
	sort.Strings(keys)
	for _, key := range keys {
		keyRune := []rune(key)[0]
		value := cmd.State.CurrentHandler().Keys[keyRune]
		cmd.State.LinesDrawn++
		if child, isKeyMap := value.(*KeyMap); isKeyMap {
			fmt.Fprintf(cmd.State.Out, "[%c] %s\n\r", keyRune, child.Name)
			continue
		}
		fmt.Fprintf(cmd.State.Out, "[%c] %s\n\r", keyRune, value)
	}
}
