package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Nerdmaster/terminal"
)

func main() {
	menuState := &MenuState{Out: os.Stdout, Err: os.Stderr, In: os.Stdin}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading current user: %s\n", err)
		return
	}
	currentDirectory, err := os.Getwd()
	if err != nil {
		currentDirectory = os.Getenv("PWD")
	}
	loadConfig := &ReloadConfig{
		Home:      currentUser.HomeDir,
		StartFrom: currentDirectory,
		State:     menuState,
	}
	menuState.DefineBuiltinCommand("quit", NewQuitCommand)
	menuState.DefineBuiltinCommand("reload", func(*MenuState) (Command, error) {
		return &ReloadConfig{
			Home:      currentUser.HomeDir,
			StartFrom: currentDirectory,
			State:     menuState,
			Verbose:   true,
		}, nil
	})

	loadConfig.Execute()
	restoreTerminal := rawTerminal()
	menuState.RestoreTerminal = restoreTerminal
	keyReader := terminal.NewKeyReader(os.Stdin)
	for !menuState.Done {
		display := &DisplayMenu{State: menuState}
		display.Execute()
		keypress, err := keyReader.ReadKeypress()
		if err != nil {
			fmt.Printf("Error reading key: %s\n", err)
			return
		}
		key := keypress.Key

		if key == terminal.KeyCtrlC {
			break
		}
		if key == terminal.KeyCtrlB || key == terminal.KeyUp || key == terminal.KeyLeft || key == terminal.KeyBackspace {
			goBack := &GoBack{State: menuState}
			goBack.Execute()
			continue
		}

		selectItem := &SelectMenuItem{
			State: menuState,
			Key:   key,
		}
		selectItem.Execute()
	}
}
