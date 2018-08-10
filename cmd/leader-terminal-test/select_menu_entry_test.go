package main_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/dhamidi/leader/cmd/leader-terminal-test"
	"github.com/stretchr/testify/assert"
)

func newTestContext(t *testing.T, root *main.KeyMap, input io.Reader) *main.Context {
	testTerminal, err := main.NewTerminalTTY()
	assert.NoError(t, err)
	return &main.Context{
		Terminal:      testTerminal.InputFrom(input),
		CurrentKeyMap: root,
	}
}

func TestSelectMenuEntry_Execute_changes_current_key_map(t *testing.T) {
	keymap := main.NewKeyMap("root")
	input := bytes.NewBufferString("a")
	keymap.Bind('a').Children().Rename("b").DefineKey('b', main.DoNothing)
	context := newTestContext(t, keymap, input)
	selectMenuEntry := main.NewSelectMenuEntry(context)
	selectMenuEntry.Execute()
	assert.Equal(t, "b", context.CurrentKeyMap.Name())
}

func TestSelectMenuEntry_Execute_runs_command_associated_with_binding(t *testing.T) {
	keymap := main.NewKeyMap("root")
	command := newMockCommand()
	input := bytes.NewBufferString("ab")
	keymap.Bind('a').Children().Rename("b").DefineKey('b', command.Execute)
	context := newTestContext(t, keymap, input)

	selectMenuEntry := main.NewSelectMenuEntry(context)
	selectMenuEntry.Execute()
	selectMenuEntry.Execute()
	assert.Equal(t, 1, command.called)

}

func TestSelectMenuEntry_Execute_gives_does_not_execute_command_on_binding_with_children(t *testing.T) {
	keymap := main.NewKeyMap("root")
	command := newMockCommand()
	input := bytes.NewBufferString("ab")
	keymap.Bind('a').Do(command.Execute).Children().DefineKey('b', main.DoNothing)
	context := newTestContext(t, keymap, input)

	selectMenuEntry := main.NewSelectMenuEntry(context)
	selectMenuEntry.Execute()
	selectMenuEntry.Execute()

	assert.Equal(t, 0, command.called)
}
