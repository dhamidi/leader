package main

// Context provides dependencies for running UI commands
type Context struct {
	Terminal      *Terminal
	CurrentKeyMap *KeyMap
	History       []*KeyMap
	KeyPath       []rune
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

// CurrentBindingIsLooping returns true if the currently selected key binding is a looping key binding.
func (ctx *Context) CurrentBindingIsLooping() bool {
	binding := ctx.CurrentKeyMap.LookupKey(ctx.KeyPath[len(ctx.KeyPath)-1])
	return binding.IsLooping()
}

// PushKey records the given key in the key path.  A key is never recorded twice in a row.
func (ctx *Context) PushKey(key rune) {
	if len(ctx.KeyPath) > 0 && ctx.KeyPath[len(ctx.KeyPath)-1] == key {
		return
	}
	ctx.KeyPath = append(ctx.KeyPath, key)
}
