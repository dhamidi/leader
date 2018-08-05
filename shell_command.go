package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

type ShellCommand struct {
	Path string
	Args []string

	Out io.Writer
	Err io.Writer
	In  io.Reader
}

func NewShellCommand(path string, args ...string) *ShellCommand {
	return &ShellCommand{
		Path: path,
		Args: args,
		Out:  os.Stdout,
		Err:  os.Stderr,
		In:   os.Stdin,
	}
}

func (cmd *ShellCommand) RedirectTo(out io.Writer, err io.Writer) *ShellCommand {
	cmd.Out = out
	cmd.Err = err
	return cmd
}

func (cmd *ShellCommand) InputFrom(in io.Reader) *ShellCommand {
	cmd.In = in
	return cmd
}

func (cmd *ShellCommand) String() string { return cmd.Path + " " + strings.Join(cmd.Args, " ") }
func (cmd *ShellCommand) Execute() {
	command := exec.Command(cmd.Path, cmd.Args...)
	command.Stdout = cmd.Out
	command.Stderr = cmd.Err
	command.Stdin = cmd.In
	command.Run()
}
