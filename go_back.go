package main

type GoBack struct {
	State *MenuState
}

func (cmd *GoBack) Execute() {
	cmd.State.PopHandler()
}
