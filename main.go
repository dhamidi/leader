package main

import (
	"fmt"
	"os"

	"github.com/Nerdmaster/terminal"
)

func main() {
	menuState := &MenuState{Out: os.Stdout, Err: os.Stderr, In: os.Stdin}

	if err := LoadLeaderRC(".leaderrc", menuState); err != nil {
		fmt.Printf("Error loading .leaderrc: %s\n", err)
		return
	}
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		fmt.Printf("Error putting terminal into raw mode: %s\n", err)
		return
	}
	defer terminal.Restore(0, oldState)
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

		selectItem := &SelectMenuItem{
			State:         menuState,
			Key:           key,
			BeforeExecute: func() { terminal.Restore(0, oldState) },
		}
		selectItem.Execute()
	}
}
