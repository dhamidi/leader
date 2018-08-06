package main

type SelectMenuItem struct {
	State *MenuState
	Key   rune
}

func (cmd *SelectMenuItem) Execute() {
	currentHandler := cmd.State.CurrentHandler()
	nextHandler, command := currentHandler.HandleKey(cmd.Key)
	if nextHandler == cmd.State.Root && command == nil {
		return
	}
	if command != nil {
		cmd.State.RestoreTerminal()
		command.Execute()
		cmd.State.Done = true
	} else {
		cmd.State.PushHandler(nextHandler)
	}
}
