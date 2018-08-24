package main

import (
	"fmt"

	"github.com/Nerdmaster/terminal"
)

// SelectMenuEntry reads a key from the terminal and changes the
// current key map to reflect this selection.
type SelectMenuEntry struct {
	*Context
}

// NewSelectMenuEntry creates a new instance of this command bound to the given context.
func NewSelectMenuEntry(ctx *Context) *SelectMenuEntry {
	return &SelectMenuEntry{
		Context: ctx,
	}
}

// Execute runs this command.
func (cmd *SelectMenuEntry) Execute() error {
	for {
		view, err := cmd.display()
		if err != nil {
			return fmt.Errorf("SelectMenuEntry: display: %s", err)
		}
		key, err := cmd.Terminal.ReadKey()
		if err != nil {
			return fmt.Errorf("SelectMenuEntry: %s", err)
		}
		if cmd.isGoBackKey(key) {
			view.Erase(cmd.Terminal)
			NewGoBack(cmd.Context).Execute()
			continue
		}
		if cmd.isExitKey(key) {
			cmd.Terminal.Restore()
			return nil
		}
		binding := cmd.CurrentKeyMap.LookupKey(key)
		cmd.PushKey(key)
		if err := view.Erase(cmd.Terminal); err != nil {
			return err
		}
		if binding.HasChildren() {
			cmd.Navigate(binding.Children())
		} else {
			if err := cmd.Terminal.Restore(); err != nil {
				return err
			}

			if binding.IsLooping() {
				if _, canLoop := cmd.Executor.(LoopingExecutor); canLoop {
					return binding.Execute()
				}
				cmd.ErrorLogger.Print(binding.Execute())
				cmd.ErrorLogger.Print(cmd.Terminal.MakeRaw())
				continue
			} else {
				return binding.Execute()
			}
		}
	}
}

// display builds the view to display before asking for input and renders it on the terminal
func (cmd *SelectMenuEntry) display() (View, error) {
	view := NewVerticalBoxView(
		cmd.displayBreadcrumbs(),
		cmd.displayMenu(),
	)

	return view, view.Render(cmd.Terminal)
}

// displayMenu builds the menu view for the currently selected keymap
func (cmd *SelectMenuEntry) displayMenu() *MenuView {
	menuEntries := []*MenuEntry{}
	for _, binding := range cmd.CurrentKeyMap.Bindings() {
		menuEntries = append(menuEntries, NewMenuEntryForKeyBinding(binding))
	}

	menu := NewMenuView(menuEntries)
	return menu
}

// displayBreadcrumbs builds the breadcrumbs view for the path to the current key map.
func (cmd *SelectMenuEntry) displayBreadcrumbs() *BreadcrumbsView {
	breadcrumbs := []string{}
	for _, keymap := range cmd.History {
		breadcrumbs = append(breadcrumbs, keymap.Name())
	}
	breadcrumbs = append(breadcrumbs, cmd.CurrentKeyMap.Name())
	breadcrumbsView := NewBreadcrumbsView(breadcrumbs)
	return breadcrumbsView
}

// isGoBackKey returns true if pressing key should go back in the menu history
func (cmd *SelectMenuEntry) isGoBackKey(key rune) bool {
	return key == terminal.KeyCtrlB ||
		key == terminal.KeyBackspace ||
		key == terminal.KeyUp ||
		key == terminal.KeyLeft
}

// isExitKey returns true if pressing key should exit the program
func (cmd *SelectMenuEntry) isExitKey(key rune) bool {
	return key == terminal.KeyCtrlC
}
