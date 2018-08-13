package main

// Context provides dependencies for running UI commands
type Context struct {
	Terminal      *Terminal
	CurrentKeyMap *KeyMap
	History       []*KeyMap
	Executor      Executor
	ErrorLogger   *ErrorLogger
}

// Navigate changes the current key map and adds the previous value to this context's history,
func (ctx *Context) Navigate(nextKeyMap *KeyMap) {
	ctx.PushHistory(ctx.CurrentKeyMap)
	ctx.CurrentKeyMap = nextKeyMap
}

// PushHistory adds a history entry to this context.
func (ctx *Context) PushHistory(keyMap *KeyMap) {
	ctx.History = append(ctx.History, keyMap)
}

// PopHistory removes the most recent history entry from this context.
func (ctx *Context) PopHistory() *KeyMap {
	if len(ctx.History) == 0 {
		return ctx.CurrentKeyMap
	}
	lastItem := ctx.History[len(ctx.History)-1]
	ctx.History = ctx.History[:len(ctx.History)-1]
	return lastItem
}
