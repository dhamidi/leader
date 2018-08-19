package main

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gobuffalo/packr"
)

// NewShellFromEnv attempts to detect the user's current shell by
// looking at the $SHELL environment variable.  If no supported shell
// is recognized, it returns a shell object configured for /bin/bash.
func NewShellFromEnv(getenv Getenv) Shell {
	shellName := filepath.Base(getenv("SHELL"))
	switch shellName {
	case "bash":
		return NewBashShell(getenv)
	case "zsh":
		return NewZSHShell(getenv)
	case "fish":
		return NewFishShell(getenv)
	default:
		return NewBashShell(getenv)
	}
}

// Getenv is a function used to access environment variables.
type Getenv func(name string) string

// POSIXShell implements Shell for /bin/bash and /bin/zsh.
type POSIXShell struct {
	name                  string
	getenv                Getenv
	commandLineBufferName string
	commandLineCursorName string
}

// NewPOSIXShell creates a new instance of bash using the provided
// function to access environment variables.
func NewPOSIXShell(name string, getenv Getenv, commandLineBufferName, commandLineCursorName string) *POSIXShell {
	return &POSIXShell{
		name:                  name,
		getenv:                getenv,
		commandLineBufferName: commandLineBufferName,
		commandLineCursorName: commandLineCursorName,
	}
}

// NewBashShell returns a POSIXShell configured for /bin/bash.
func NewBashShell(getenv Getenv) *POSIXShell {
	return NewPOSIXShell("bash", getenv, "READLINE_LINE", "READLINE_POINT")
}

// NewZSHShell returns a POSIXShell configured for /bin/zsh.
func NewZSHShell(getenv Getenv) *POSIXShell {
	return NewPOSIXShell("zsh", getenv, "BUFFER", "CURSOR")
}

// Commandline returns the current line of input and cursor position
// in bash by looking at the READLINE_INPUT and READLINE_POINT
// variables, as per bash's documentation.
func (shell *POSIXShell) Commandline() (string, int) {
	line := shell.getenv(shell.commandLineBufferName)
	cursorString := shell.getenv(shell.commandLineCursorName)
	cursor, err := strconv.Atoi(cursorString)
	if err != nil {
		return line, 0
	}

	return line, cursor
}

// EvalNext returns a shell command for running cmd followed by
// invoking leader again at the specified path.
//
// Example: shell.EvalNext("date", "/") => `date; eval "$(leader print @/)"`
func (shell *POSIXShell) EvalNext(cmd string, path []rune) string {
	return fmt.Sprintf("%s; eval \"$(leader print @%s)\"", cmd, string(path))
}

// Init returns initialization code for installing leader in bash.
func (shell *POSIXShell) Init() string {
	return packr.NewBox("assets").String(fmt.Sprintf("leader.%s.sh", shell.name))
}

// FishShell implements Shell for /usr/bin/fish.
type FishShell struct {
	getenv Getenv
}

// NewFishShell returns a new FishShell using the provided command runner.
func NewFishShell(getenv Getenv) *FishShell {
	return &FishShell{
		getenv: getenv,
	}
}

// Commandline returns the current input line and cursor position by
// invoking inspecting FISH_INPUT and FISH_POINT.  Both of these
// variables are expected to be set by the wrapper function which is
// bound to the leader key in leader.fish.sh.
func (shell *FishShell) Commandline() (string, int) {
	line := shell.getenv("FISH_INPUT")
	cursor, err := strconv.Atoi(shell.getenv("FISH_POINT"))
	if err != nil {
		return line, 0
	}

	return line, cursor

}

// EvalNext returns a shell command for running cmd followed by
// invoking leader again at the specified path.
//
// Example: shell.EvalNext("date", "/") => `date; eval "$(leader print @/)"`
func (shell *FishShell) EvalNext(cmd string, path []rune) string {
	return fmt.Sprintf("%s; eval (leader print @%s)", cmd, string(path))
}

// Init returns initialization code for installing leader in fish.
func (shell *FishShell) Init() string {
	return packr.NewBox("assets").String("leader.fish.sh")
}
