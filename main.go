package main

import (
	"fmt"
	"os"

	"github.com/Nerdmaster/terminal"
)

func main() {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		fmt.Printf("Error putting terminal into raw mode: %s\n", err)
		return
	}
	defer terminal.Restore(0, oldState)

	done := false
	keyMap := &KeyMap{
		Name: "global",
		Keys: map[rune]interface{}{
			'q': &QuitCommand{DoneVar: &done},
			'g': &KeyMap{
				Name: "go",
				Keys: map[rune]interface{}{
					'b': NewShellCommand("go", "build", "."),
					't': NewShellCommand("go", "test", "."),
				},
			},
		},
	}

	keyReader := terminal.NewKeyReader(os.Stdin)
	menuState := &MenuState{Out: os.Stdout, KeyMap: keyMap}
	for !done {
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
			AfterExecute:  func() { done = true },
		}
		selectItem.Execute()
	}
}
