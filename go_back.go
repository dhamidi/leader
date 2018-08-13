package main

// GoBack goes back to the previous key map in the given context.
type GoBack struct {
	*Context
}

// NewGoBack creates a new instance of this command operating on context.
func NewGoBack(ctx *Context) *GoBack {
	return &GoBack{
		Context: ctx,
	}
}

// Execute runs this command.
func (cmd *GoBack) Execute() error {
	cmd.CurrentKeyMap = cmd.PopHistory()
	return nil
}
