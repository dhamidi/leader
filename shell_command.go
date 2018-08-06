package main

import (
	"fmt"
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

		Out: os.Stdout,
		Err: os.Stderr,
		In:  os.Stdin,
	}
}

func (cmd *ShellCommand) InputFrom(in io.Reader) *ShellCommand {
	cmd.In = in
	return cmd
}
func (cmd *ShellCommand) RedirectTo(out io.Writer, err io.Writer) *ShellCommand {
	cmd.Out = out
	cmd.Err = err
	return cmd
}
func (cmd *ShellCommand) String() string { return cmd.Path + " " + strings.Join(cmd.Args, " ") }
func (cmd *ShellCommand) Execute() {
	args := append([]string{cmd.Path}, cmd.Args...)
	for i, arg := range args {
		args[i] = fmt.Sprintf(`"%s"`, strings.Replace(arg, `"`, `\"`, 0))
	}
	asScript := strings.Join(args, " ")
	command := exec.Command("bash", "-c", asScript)
	command.Stdout = cmd.Out
	command.Stderr = cmd.Err
	command.Stdin = cmd.In
	command.Run()
}
