package main_test

import (
	"bytes"
	"testing"

	"github.com/dhamidi/leader"
)

type testbed struct {
	out *bytes.Buffer
	err *bytes.Buffer
	in  *bytes.Buffer
	t   *testing.T
}

func newTestbed(t *testing.T) *testbed {
	result := &testbed{
		out: bytes.NewBufferString(""),
		err: bytes.NewBufferString(""),
		in:  bytes.NewBufferString(""),
		t:   t,
	}
	return result
}

type testcase struct {
	t     *testing.T
	bed   *testbed
	state *main.MenuState
}

func newTestCase(t *testing.T) *testcase {
	result := &testcase{
		t:   t,
		bed: newTestbed(t),
	}
	result.state = newMenuState(result.bed)

	return result
}

func newMenuState(bed *testbed) *main.MenuState {
	menuState := &main.MenuState{
		Out:             bed.out,
		Err:             bed.err,
		In:              bed.in,
		RestoreTerminal: func() {},
		Done:            false,
	}

	menuState.Root = &main.KeyMap{
		Name: "global",
		Keys: map[rune]interface{}{
			'd': &dummyCommand{},
			'a': &main.KeyMap{
				Name: "a",
				Keys: map[rune]interface{}{
					'b': &main.QuitCommand{State: menuState},
				},
			},
		},
	}

	return menuState
}

type dummyCommand struct {
	executed int
}

func (cmd *dummyCommand) Execute() {
	cmd.executed++
}
