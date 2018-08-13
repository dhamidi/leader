package main_test

import (
	"fmt"
	"testing"

	"github.com/dhamidi/leader/cmd/leader-terminal-test"
	"github.com/stretchr/testify/assert"
)

type mockCommand struct{ called int }

func newMockCommand() *mockCommand { return &mockCommand{called: 0} }
func (m *mockCommand) Execute() error {
	m.called++
	return nil
}

func TestKeyMap_DefineKey_creates_an_entry_at_the_current_map_level(t *testing.T) {
	keymap := main.NewKeyMap("root")
	keymap.DefineKey('a', main.DoNothing)
	assert.NotEqual(t, main.UnboundKey, keymap.LookupKey('a'))
}

func TestKeyMap_DefineKey_associates_a_binding_with_the_provided_command(t *testing.T) {
	keymap := main.NewKeyMap("root")
	keymap.DefineKey('a', main.FailWhenExecuted(fmt.Errorf("done")))
	assert.Equal(t, fmt.Errorf("done"), keymap.LookupKey('a').Execute())
}

func TestKeyMap_nested_bindings_are_supported(t *testing.T) {
	keymap := main.NewKeyMap("root")
	command := newMockCommand()
	keymap.Bind('a').Children().Bind('b').Do(command.Execute)
	keymap.LookupKey('a').Children().LookupKey('b').Execute()
	assert.Equal(t, 1, command.called, "command has not been called")
}
