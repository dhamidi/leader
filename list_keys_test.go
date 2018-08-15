package main_test

import (
	"bytes"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestListKeys_Execute_prints_all_key_bindings(t *testing.T) {
	keymap := main.NewKeyMap("root")
	keymapA := keymap.Bind('a').Children()
	keymapA.Bind('a').Describe("run 'a a'")
	keymapA.Bind('b').Describe("run 'a b'")
	keymapB := keymap.Bind('b').Children()
	keymapBB := keymapB.Bind('b').Children()
	keymapBB.Bind('b').Describe("run 'b b b'")

	output := bytes.NewBufferString("")
	context := newTestContext(t, keymap, bytes.NewBufferString(""), output)

	listKeys := main.NewListKeys(context)
	listKeys.Execute()

	expectedOutput := `a a: run 'a a'
a b: run 'a b'
b b b: run 'b b b'`
	assert.Contains(t, output.String(), expectedOutput)
}
