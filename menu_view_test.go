package main_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestMenuView_Render_it_shows_each_menu_entry_on_a_separate_line(t *testing.T) {
	out := bytes.NewBufferString("")
	view := main.NewMenuView([]*main.MenuEntry{
		{Key: 'a', Label: "command a"},
		{Key: 'b', Label: "command b"},
	})

	assert.NoError(t, view.Render(out))
	lines := strings.Split(out.String(), "\n\r")
	assert.ElementsMatch(t, []string{
		"[a] command a",
		"[b] command b",
		"",
	}, lines)
}

func TestMenuView_Erase_erases_one_line_per_entry(t *testing.T) {
	out := bytes.NewBufferString("")
	view := main.NewMenuView([]*main.MenuEntry{
		{Key: 'a', Label: "command a"},
		{Key: 'b', Label: "command b"},
	})

	assert.NoError(t, view.Erase(out))
	assert.Equal(t, "\033[A\033[2K\033[A\033[2K", out.String())
}
