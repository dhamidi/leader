package main

import (
	"fmt"

	"github.com/Nerdmaster/terminal"
)

var keyText = map[rune]string{
	terminal.KeyCtrlA:        "Ctrl+A",
	terminal.KeyCtrlB:        "Ctrl+B",
	terminal.KeyCtrlC:        "Ctrl+C",
	terminal.KeyCtrlD:        "Ctrl+D",
	terminal.KeyCtrlE:        "Ctrl+E",
	terminal.KeyCtrlF:        "Ctrl+F",
	terminal.KeyCtrlG:        "Ctrl+G",
	terminal.KeyCtrlH:        "Ctrl+H",
	terminal.KeyCtrlI:        "Ctrl+I",
	terminal.KeyCtrlJ:        "Ctrl+J",
	terminal.KeyCtrlK:        "Ctrl+K",
	terminal.KeyCtrlL:        "Ctrl+L",
	terminal.KeyCtrlN:        "Ctrl+N",
	terminal.KeyCtrlO:        "Ctrl+O",
	terminal.KeyCtrlP:        "Ctrl+P",
	terminal.KeyCtrlQ:        "Ctrl+Q",
	terminal.KeyCtrlR:        "Ctrl+R",
	terminal.KeyCtrlS:        "Ctrl+S",
	terminal.KeyCtrlT:        "Ctrl+T",
	terminal.KeyCtrlU:        "Ctrl+U",
	terminal.KeyCtrlV:        "Ctrl+V",
	terminal.KeyCtrlW:        "Ctrl+W",
	terminal.KeyCtrlX:        "Ctrl+X",
	terminal.KeyCtrlY:        "Ctrl+Y",
	terminal.KeyCtrlZ:        "Ctrl+Z",
	terminal.KeyEscape:       "Escape",
	terminal.KeyLeftBracket:  "[",
	terminal.KeyRightBracket: "]",
	terminal.KeyEnter:        "Enter",
	terminal.KeyBackspace:    "Backspace",
	terminal.KeyUnknown:      "Unknown",
	terminal.KeyUp:           "Up",
	terminal.KeyDown:         "Down",
	terminal.KeyLeft:         "Left",
	terminal.KeyRight:        "Right",
	terminal.KeyHome:         "Home",
	terminal.KeyEnd:          "End",
	terminal.KeyPasteStart:   "PasteStart",
	terminal.KeyPasteEnd:     "PasteEnd",
	terminal.KeyInsert:       "Insert",
	terminal.KeyDelete:       "Delete",
	terminal.KeyPgUp:         "PgUp",
	terminal.KeyPgDn:         "PgDn",
	terminal.KeyPause:        "Pause",
	terminal.KeyF1:           "F1",
	terminal.KeyF2:           "F2",
	terminal.KeyF3:           "F3",
	terminal.KeyF4:           "F4",
	terminal.KeyF5:           "F5",
	terminal.KeyF6:           "F6",
	terminal.KeyF7:           "F7",
	terminal.KeyF8:           "F8",
	terminal.KeyF9:           "F9",
	terminal.KeyF10:          "F10",
	terminal.KeyF11:          "F11",
	terminal.KeyF12:          "F12",
}

func keyFromString(keyDescription string) rune {
	for key, description := range keyText {
		if description == keyDescription {
			return key
		}
	}

	return ([]rune(keyDescription))[0]
}

func keyToString(key rune) string {
	text, found := keyText[key]
	if found {
		return text
	}

	return fmt.Sprintf("%c", key)
}
