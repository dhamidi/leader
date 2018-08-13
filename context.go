package main

// Context provides dependencies for running UI commands
type Context struct {
	Terminal      *Terminal
	CurrentKeyMap *KeyMap
	Executor      Executor
	ErrorLogger   *ErrorLogger
}
