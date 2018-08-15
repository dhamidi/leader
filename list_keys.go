package main

import (
	"fmt"
	"strings"
)

// ListKeys prints all key bindings on the context's terminal.
type ListKeys struct {
	*Context
}

// NewListKeys creates a new command instance operating on the provided context.
func NewListKeys(ctx *Context) *ListKeys {
	return &ListKeys{
		Context: ctx,
	}
}

// Execute runs this command.
func (cmd *ListKeys) Execute() error {
	var walkKeyMap func(keymap *KeyMap, path []rune, f func(path []rune, b *KeyBinding))
	walkKeyMap = func(keymap *KeyMap, path []rune, f func(path []rune, b *KeyBinding)) {
		for _, keyBinding := range keymap.Bindings() {
			childPath := make([]rune, len(path))
			copy(childPath, path)
			childPath = append(childPath, keyBinding.Key())
			if keyBinding.HasChildren() {
				walkKeyMap(keyBinding.Children(), childPath, f)
			} else {
				f(childPath, keyBinding)
			}
		}
	}

	printBinding := func(path []rune, b *KeyBinding) {
		pathString := []string{}
		for _, r := range path {
			pathString = append(pathString, fmt.Sprintf("%c", r))
		}
		fmt.Fprintf(cmd.Terminal, "%s: %s\n", strings.Join(pathString, " "), b.Description())
	}

	walkKeyMap(cmd.CurrentKeyMap, []rune{}, printBinding)
	return nil
}
