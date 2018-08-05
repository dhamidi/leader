package main

import (
	"os"
	"os/exec"
	"strings"
)

type ShellCommand struct {
	Path string
	Args []string
}

func NewShellCommand(path string, args ...string) *ShellCommand {
	return &ShellCommand{
		Path: path,
		Args: args,
	}
}

func (cmd *ShellCommand) String() string { return cmd.Path + " " + strings.Join(cmd.Args, " ") }
func (cmd *ShellCommand) Execute() {
	command := exec.Command(cmd.Path, cmd.Args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin
	command.Run()
}
