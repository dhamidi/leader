package main_test

import (
	"bytes"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestGoBack_Execute_selects_the_previously_selected_keymap(t *testing.T) {
	keymap := main.NewKeyMap("root")
	keymap.
		Bind('a').Children().Rename("a").
		Bind('b').Children().Rename("b").
		Bind('c').Describe("c")
	input := bytes.NewBufferString("ab")
	context := newTestContext(t, keymap, input)
	selectMenuEntry := main.NewSelectMenuEntry(context)
	selectMenuEntry.Execute()
	goBack := main.NewGoBack(context)
	assert.NoError(t, goBack.Execute())

	assert.Equal(t, "a", context.CurrentKeyMap.Name())
}
