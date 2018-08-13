package main

import (
	"io"
	"os"
	"os/exec"
)

// ShellExecutor implements Executor by running commands in a shell.
type ShellExecutor struct {
	shell string
	args  []string

	in  io.Reader
	out io.Writer
	err io.Writer
}

// NewShellExecutor creates a new executor for running commands in shell.
//
// The value of shell should either be full path to a shell executable
// (e.g. /bin/bash) or just the name of executable (e.g. bash).
//
// The value of args should be the command line flags accepted by
// shell that are necessary to run shell commands.
func NewShellExecutor(shell string, args ...string) *ShellExecutor {
	return &ShellExecutor{
		shell: shell,
		args:  args,
	}
}

// Attach sets the input and output channels for this executor to the given file.
func (e *ShellExecutor) Attach(f *os.File) *ShellExecutor {
	e.in = f
	e.out = f
	e.err = f

	return e
}

// RunCommand implements Executor by invoking shellCommand using `bash -c`
func (e *ShellExecutor) RunCommand(shellCommand string) error {
	args := []string{}
	for _, arg := range e.args {
		args = append(args, arg)
	}
	args = append(args, shellCommand)
	command := exec.Command(e.shell, args...)
	command.Stdin = e.in
	command.Stdout = e.out
	command.Stderr = e.err
	return command.Run()
}
