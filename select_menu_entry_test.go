package main_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/Nerdmaster/terminal"
	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

type testTerminal struct {
	In        io.Reader
	Out       io.Writer
	KeyReader *terminal.KeyReader
}

func newTestTerminal() *testTerminal {
	return &testTerminal{
		In:  bytes.NewBufferString(""),
		Out: bytes.NewBufferString(""),
	}
}

func (term *testTerminal) MakeRaw() error              { return nil }
func (term *testTerminal) Restore() error              { return nil }
func (term *testTerminal) Write(p []byte) (int, error) { return term.Out.Write(p) }
func (term *testTerminal) OutputTo(out io.Writer) *testTerminal {
	term.Out = out
	return term
}
func (term *testTerminal) InputFrom(in io.Reader) *testTerminal {
	term.In = in
	term.KeyReader = terminal.NewKeyReader(in)
	return term
}

func (term *testTerminal) ReadKey() (rune, error) {
	key, err := term.KeyReader.ReadKeypress()
	if err != nil {
		return terminal.KeyCtrlC, nil
	}
	return key.Key, nil
}

func newTestContext(t *testing.T, root *main.KeyMap, input io.Reader, output io.Writer) *main.Context {
	testTerminal := newTestTerminal().InputFrom(input)
	if output != nil {
		testTerminal.OutputTo(output)
	}

	return &main.Context{
		Terminal:      testTerminal,
		CurrentKeyMap: root,
	}
}

func TestSelectMenuEntry_Execute_changes_current_key_map(t *testing.T) {
	keymap := main.NewKeyMap("root")
	input := bytes.NewBufferString("a")
	keymap.Bind('a').Children().Rename("b").DefineKey('b', main.DoNothing)
	context := newTestContext(t, keymap, input, nil)
	selectMenuEntry := main.NewSelectMenuEntry(context)
	selectMenuEntry.Execute()
	assert.Equal(t, "b", context.CurrentKeyMap.Name())
}

func TestSelectMenuEntry_Execute_runs_command_associated_with_binding(t *testing.T) {
	keymap := main.NewKeyMap("root")
	command := newMockCommand()
	input := bytes.NewBufferString("ab")
	keymap.Bind('a').Children().Rename("b").DefineKey('b', command.Execute)
	context := newTestContext(t, keymap, input, nil)

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
	context := newTestContext(t, keymap, input, nil)

	selectMenuEntry := main.NewSelectMenuEntry(context)
	selectMenuEntry.Execute()
	selectMenuEntry.Execute()

	assert.Equal(t, 0, command.called)
}

func TestSelectMenuEntry_Execute_displays_breadscrumbs_for_the_current_path(t *testing.T) {
	keymap := main.NewKeyMap("root")
	input := bytes.NewBufferString("a")
	output := bytes.NewBufferString("")
	keymap.Bind('a').Children().Rename("a").DefineKey('b', main.DoNothing)
	context := newTestContext(t, keymap, input, output)
	selectMenuEntry := main.NewSelectMenuEntry(context)

	selectMenuEntry.Execute()
	assert.Contains(t, output.String(), "root > a")
}

func TestSelectMenuEntry_Execute_displays_the_current_keymap_as_a_menu(t *testing.T) {
	keymap := main.NewKeyMap("root")
	keymap.Bind('a').Do(main.DoNothing).Describe("do a")
	keymap.Bind('b').Do(main.DoNothing).Describe("do b")
	input := bytes.NewBufferString("")
	output := bytes.NewBufferString("")
	context := newTestContext(t, keymap, input, output)

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
	output := bytes.NewBufferString("")
	context := newTestContext(t, keymap, input, output)
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

func TestSelectMenuEntry_Execute_erases_the_current_menu_before_running_a_command(t *testing.T) {
	keymap := main.NewKeyMap("root")
	input := bytes.NewBufferString("ab")
	keymap.Bind('a').Children().Rename("b").DefineKey('b', main.DoNothing)
	output := bytes.NewBufferString("")
	context := newTestContext(t, keymap, input, output)
	selectMenuEntry := main.NewSelectMenuEntry(context)

	selectMenuEntry.Execute()

	expectedViews := []main.View{
		main.NewBreadcrumbsView([]string{"root"}),
		main.NewMenuView([]*main.MenuEntry{
			{Key: 'a', Label: "b"},
		}),
		main.NewBreadcrumbsView([]string{"root", "b"}),
		main.NewMenuView([]*main.MenuEntry{
			{Key: 'b', Label: ""},
		}),
	}
	outputBuffer := bytes.NewBufferString("")
	expectedViews[0].Render(outputBuffer)
	expectedViews[1].Render(outputBuffer)
	expectedViews[0].Erase(outputBuffer)
	expectedViews[1].Erase(outputBuffer)
	expectedViews[2].Render(outputBuffer)
	expectedViews[3].Render(outputBuffer)
	expectedViews[2].Erase(outputBuffer)
	expectedViews[3].Erase(outputBuffer)

	assert.True(t, bytes.Contains(output.Bytes(), outputBuffer.Bytes()),
		"output %q does not contain instructions %q", output, outputBuffer)
}

func TestSelectMenuEntry_Execute_keeps_executing_looping_keys_repeatedly(t *testing.T) {
	command := newMockCommand()
	input := bytes.NewBufferString("aaaa")
	keymap := main.NewKeyMap("root")
	keymap.Bind('a').Do(command.Execute).SetLooping(true)
	context := newTestContext(t, keymap, input, nil)
	selectMenuEntry := main.NewSelectMenuEntry(context)

	selectMenuEntry.Execute()

	assert.Equal(t, 4, command.called)
}
