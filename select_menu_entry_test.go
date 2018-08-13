package main_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func newTestContext(t *testing.T, root *main.KeyMap, input io.Reader) *main.Context {
	testTerminal, err := main.NewTerminalTTY()
	assert.NoError(t, err)
	return &main.Context{
		Terminal:      testTerminal.InputFrom(input).OutputTo(bytes.NewBufferString("")),
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

func TestSelectMenuEntry_Execute_displays_the_current_keymap_as_a_menu(t *testing.T) {
	keymap := main.NewKeyMap("root")
	keymap.Bind('a').Do(main.DoNothing).Describe("do a")
	keymap.Bind('b').Do(main.DoNothing).Describe("do b")
	input := bytes.NewBufferString("")
	output := bytes.NewBufferString("")
	context := newTestContext(t, keymap, input)
	context.Terminal.OutputTo(output)

	selectMenuEntry := main.NewSelectMenuEntry(context)
	selectMenuEntry.Execute()

	expectedMenu := main.NewMenuView([]*main.MenuEntry{
		{Key: 'a', Label: "do a"},
		{Key: 'b', Label: "do b"},
	})
	expectedOutput := main.MustRenderViewToString(expectedMenu)

	assert.Contains(t, output.String(), expectedOutput)
}

func TestSelectMenuEntry_Execute_erases_the_current_menu_before_selecting_a_child_menu(t *testing.T) {
	keymap := main.NewKeyMap("root")
	input := bytes.NewBufferString("a")
	keymap.Bind('a').Children().Rename("b").DefineKey('b', main.DoNothing)
	context := newTestContext(t, keymap, input)
	output := bytes.NewBufferString("")
	context.Terminal.OutputTo(output)
	selectMenuEntry := main.NewSelectMenuEntry(context)

	selectMenuEntry.Execute()

	expectedMenu := main.NewMenuView([]*main.MenuEntry{
		{Key: 'a', Label: "do a"},
	})
	eraseMenuBuffer := bytes.NewBufferString("")
	expectedMenu.Erase(eraseMenuBuffer)

	assert.True(t, bytes.Contains(output.Bytes(), eraseMenuBuffer.Bytes()),
		"output %q does not contain instructions %q", output, eraseMenuBuffer)
}
