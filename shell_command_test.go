package main_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
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

func shellLevel(t *testing.T) int {
	t.Helper()
	asText := os.Getenv("SHLVL")
	asInt, err := strconv.Atoi(asText)
	if err != nil {
		t.Fatalf("Failed to read shell level from $SHLVL=%s: %s", asText, err)
	}
	return asInt
}

func TestShellCommand_Execute_runs_command_in_bash(t *testing.T) {
	currentShellLevel := shellLevel(t)
	command := main.NewShellCommand("printf", "%s %s", "$SHELL", "$SHLVL")
	testbed := newTestbed(t)
	command.RedirectTo(testbed.out, testbed.err).InputFrom(testbed.in).Execute()
	shell, actualShellLevel := "", 0
	if _, err := fmt.Sscanf(testbed.out.String(), "%s %d", &shell, &actualShellLevel); err != nil {
		t.Fatalf("error parsing output %q: %s", testbed.out.String(), err)
	}
	assert.Equal(t, "bash", filepath.Base(shell))
	assert.Equal(t, currentShellLevel+1, actualShellLevel)
}

func TestShellCommand_Execute_runs_command_in_the_shell_configured_by_the_user(t *testing.T) {
	os.Setenv("SHELL", "echo")
	testbed := newTestbed(t)
	command := main.NewShellCommand("true")
	command.RedirectTo(testbed.out, testbed.err).InputFrom(testbed.in).Execute()

	assert.Equal(t, `-c "true"`, strings.TrimSpace(testbed.out.String()))
}
