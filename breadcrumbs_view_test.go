package main_test

import (
	"bytes"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestBreadcrumbsView_Render_shows_each_breadcrumb_separated_with_a_greater_then_sign(t *testing.T) {
	view := main.NewBreadcrumbsView([]string{"root", "ruby", "rake"})
	output := bytes.NewBufferString("")
	assert.NoError(t, view.Render(output))
	assert.Contains(t, output.String(), "root > ruby > rake")
}

func TestBreadcrumbsView_Erase_erases_one_line(t *testing.T) {
	view := main.NewBreadcrumbsView([]string{"root", "ruby", "rake"})
	output := bytes.NewBufferString("")
	assert.NoError(t, view.Erase(output))
	assert.True(t, bytes.Contains(output.Bytes(), []byte("\033[A\033[2K")), "expected %q to contain instructions %q", output, "\033[A\033[2K")
}
