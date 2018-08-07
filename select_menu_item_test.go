package main_test

import (
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestSelectMenuItem_Execute_selects_the_next_key_map_based_on_the_pressed_key(t *testing.T) {
	testcase := newTestCase(t)
	selectMenuItem := &main.SelectMenuItem{
		Key:   rune('a'),
		State: testcase.state,
	}

	selectMenuItem.Execute()

	assert.Equal(t, "a", testcase.state.CurrentPath())
}

func TestSelectMenuItem_Execute_does_not_change_the_map_if_the_pressed_key_is_not_bound(t *testing.T) {
	testcase := newTestCase(t)
	selectMenuItem := &main.SelectMenuItem{
		Key:   rune('z'),
		State: testcase.state,
	}
	selectMenuItem.Execute()

	assert.Equal(t, "", testcase.state.CurrentPath())
}

func TestSelectMenuItem_Execute_runs_the_associated_command_if_key_is_not_mapped_to_a_keymap(t *testing.T) {
	testcase := newTestCase(t)
	selectSubmenu := &main.SelectMenuItem{
		Key:   rune('a'),
		State: testcase.state,
	}
	selectSubmenu.Execute()
	quit := &main.SelectMenuItem{
		Key:   rune('b'),
		State: testcase.state,
	}
	quit.Execute()

	assert.Equal(t, true, testcase.state.Done, "state.done")
}

func TestSelectMenuItem_Execute_restores_the_terminal_state_before_running_a_command(t *testing.T) {
	testcase := newTestCase(t)
	runDummyCommand := &main.SelectMenuItem{
		Key:   rune('d'),
		State: testcase.state,
	}
	terminalRestored := false
	testcase.state.RestoreTerminal = func() {
		terminalRestored = true
	}

	runDummyCommand.Execute()

	assert.Equal(t, true, terminalRestored, "terminalRestored")

}
