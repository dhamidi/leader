package main

import (
	"fmt"
	"io"
)

// MenuView is a view of the current key map.
type MenuView struct {
	// Entries are the menu entries that should be displayed by
	// this view.
	Entries []*MenuEntry
}

// NewMenuView constructs a menu view for rendering Entries.
func NewMenuView(entries []*MenuEntry) *MenuView {
	return &MenuView{
		Entries: entries,
	}
}

// Render displays the menu on out.
func (v *MenuView) Render(out io.Writer) error {
	for _, entry := range v.Entries {
		_, err := fmt.Fprintf(out, "[%c] %s\n\r", entry.Key, entry.Label)
		if err != nil {
			return err
		}
	}

	return nil
}

// MenuEntry represents an entry in a menu that should be displayed in a terminal.
type MenuEntry struct {
	// Key that needs to be pressed to select this menu entry.
	Key rune

	// Label that is shown to the user, describing the menu entry
	Label string
}

// NewMenuEntryForKeyBinding returns a new menu entry for a given key binding.
func NewMenuEntryForKeyBinding(binding *KeyBinding) *MenuEntry {
	return &MenuEntry{
		Key:   binding.key,
		Label: binding.description,
	}
}
