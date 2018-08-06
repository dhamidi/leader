package main

import (
	"io"
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
