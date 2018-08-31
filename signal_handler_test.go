package main_test

import (
	"syscall"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestSignalHandler_does_not_exit_on_SIGWINCH(t *testing.T) {
	exited := false
	exitFn := func(int) { exited = true }
	main.SignalHandler(exitFn)(syscall.SIGWINCH)
	assert.Equal(t, false, exited)
}

func TestSignalHandler_exits_with_code_0_on_any_other_signal(t *testing.T) {
	exited := -1
	exitFn := func(exitCode int) { exited = exitCode }
	main.SignalHandler(exitFn)(syscall.SIGINT)
	assert.Equal(t, 0, exited)
}
