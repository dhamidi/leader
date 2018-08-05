package main

type QuitCommand struct {
	DoneVar *bool
}

func (cmd *QuitCommand) String() string { return "quit" }

func (cmd *QuitCommand) Execute() {
	*cmd.DoneVar = true
}
