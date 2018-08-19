package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestBash_Commandline_extracts_info_from_environment(t *testing.T) {
	getenv := func(name string) string {
		switch name {
		case "READLINE_LINE":
			return "line"
		case "READLINE_POINT":
			return "1"
		default:
			return ""
		}
	}
	bash := main.NewBashShell(getenv)

	line, cursor := bash.Commandline()
	assert.Equal(t, "line", line)
	assert.Equal(t, 1, cursor)
}

func TestBash_EvalNext_emits_eval_command_for_invoking_leader(t *testing.T) {
	bash := main.NewBashShell(os.Getenv)
	evalNext := bash.EvalNext("date", []rune("/"))
	assert.Equal(t, `date; eval "$(leader print @/)"`, evalNext)
}

func TestBash_Init_returns_shell_script_for_bash_initialization(t *testing.T) {
	bash := main.NewBashShell(os.Getenv)
	initCode, err := ioutil.ReadFile("./assets/leader.bash.sh")
	assert.NoError(t, err)
	assert.Equal(t, string(initCode), bash.Init())
}

func TestZSH_Commandline_extracts_info_from_environment(t *testing.T) {
	getenv := func(name string) string {
		switch name {
		case "BUFFER":
			return "line"
		case "CURSOR":
			return "1"
		default:
			return ""
		}
	}
	zsh := main.NewZSHShell(getenv)

	line, cursor := zsh.Commandline()
	assert.Equal(t, "line", line)
	assert.Equal(t, 1, cursor)
}

func TestZSH_EvalNext_emits_eval_command_for_invoking_leader(t *testing.T) {
	zsh := main.NewZSHShell(os.Getenv)
	evalNext := zsh.EvalNext("date", []rune("/"))
	assert.Equal(t, `date; eval "$(leader print @/)"`, evalNext)
}

func TestZSH_Init_returns_shell_script_for_bash_initialization(t *testing.T) {
	zsh := main.NewZSHShell(os.Getenv)
	initCode, err := ioutil.ReadFile("./assets/leader.zsh.sh")
	assert.NoError(t, err)
	assert.Equal(t, string(initCode), zsh.Init())
}

func TestFish_Commandline_extracts_info_by_running_commandline(t *testing.T) {
	getenv := func(name string) string {
		switch name {
		case "FISH_INPUT":
			return "line"
		case "FISH_POINT":
			return "1"
		default:
			return ""
		}
	}
	fish := main.NewFishShell(getenv)

	line, cursor := fish.Commandline()
	assert.Equal(t, "line", line)
	assert.Equal(t, 1, cursor)
}

func TestFish_EvalNext_emits_eval_command_for_invoking_leader(t *testing.T) {
	fish := main.NewFishShell(os.Getenv)
	evalNext := fish.EvalNext("date", []rune("/"))
	assert.Equal(t, `date; eval (leader print @/)`, evalNext)
}

func TestFish_Init_returns_shell_script_for_bash_initialization(t *testing.T) {
	fish := main.NewFishShell(os.Getenv)
	initCode, err := ioutil.ReadFile("./assets/leader.fish.sh")
	assert.NoError(t, err)
	assert.Equal(t, string(initCode), fish.Init())
}
