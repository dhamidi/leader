package main

type QuitCommand struct {
	State *MenuState
}

func NewQuitCommand(state *MenuState) (Command, error) { return &QuitCommand{State: state}, nil }
func (cmd *QuitCommand) String() string                { return "quit" }

func (cmd *QuitCommand) Execute() {
	cmd.State.Done = true
}
