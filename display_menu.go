package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type MenuState struct {
	Done            bool
	Out             io.Writer
	Err             io.Writer
	In              io.Reader
	LinesDrawn      int
	Root            *KeyMap
	Path            []*KeyMap
	BuiltinCommands map[string]CommandFn
	RestoreTerminal func()
}

func (m *MenuState) CurrentPath() string {
	names := []string{}
	for _, keyMap := range m.Path {
		names = append(names, keyMap.Name)
	}

	return strings.Join(names, " > ")
}
func (m *MenuState) PopHandler() {
	if len(m.Path) == 0 {
		return
	}

	m.Path = m.Path[:len(m.Path)-1]
}
func (m *MenuState) PushHandler(keyMap *KeyMap) {
	if m.CurrentHandler() == keyMap {
		return
	}
	m.Path = append(m.Path, keyMap)
}
func (m *MenuState) CurrentHandler() *KeyMap {
	if len(m.Path) == 0 {
		return m.Root
	}
	return m.Path[len(m.Path)-1]
}

func (m *MenuState) DefineBuiltinCommand(name string, fn CommandFn) {
	if m.BuiltinCommands == nil {
		m.BuiltinCommands = map[string]CommandFn{}
	}
	m.BuiltinCommands[name] = fn
}

type QuitCommand struct {
	State *MenuState
}

func NewQuitCommand(state *MenuState) (Command, error) { return &QuitCommand{State: state}, nil }
func (cmd *QuitCommand) String() string                { return "quit" }

func (cmd *QuitCommand) Execute() {
	cmd.State.Done = true
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

type SelectMenuItem struct {
	State *MenuState
	Key   rune
}

func (cmd *SelectMenuItem) Execute() {
	currentHandler := cmd.State.CurrentHandler()
	nextHandler, command := currentHandler.HandleKey(cmd.Key)
	if nextHandler == cmd.State.Root && command == nil {
		return
	}
	if command != nil {
		cmd.State.RestoreTerminal()
		command.Execute()
		cmd.State.Done = true
	} else {
		cmd.State.PushHandler(nextHandler)
	}
}

type GoBack struct {
	State *MenuState
}

func (cmd *GoBack) Execute() {
	cmd.State.PopHandler()
}
