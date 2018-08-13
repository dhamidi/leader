package main

import (
	"fmt"
	"io"
)

// PrintingExecutor implements Executor by printing commands instead of running them.
type PrintingExecutor struct {
	*Context
	out io.Writer
}

// NewPrintingExecutor creates an executor that prints commands on out.
func NewPrintingExecutor(context *Context, out io.Writer) *PrintingExecutor {
	return &PrintingExecutor{
		Context: context,
		out:     out,
	}
}

// RunCommand prints the command on cmd.
func (e *PrintingExecutor) RunCommand(shellCommand string) error {
	fmt.Fprintf(e.out, "%s\n", shellCommand)
	return nil
}

// RunLoopingCommand runs a looping command by appending an `eval` for
// the current leader state to the shell command in question.
func (e *PrintingExecutor) RunLoopingCommand(shellCommand string) error {
	fmt.Fprintf(e.out, "%s; eval \"$(leader print @%s)\"", shellCommand, string(e.KeyPath[:len(e.KeyPath)-1]))
	return nil
}
