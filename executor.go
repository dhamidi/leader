package main

// Executor defines the interface for starting external one-off
// processes, i.e. executables that are spawned by this program.
type Executor interface {
	// RunCommand runs a shell command without performing any kind of escaping.
	//
	// The terminal's stdout, stderr and stdin are attached to the
	// command that is spawned by RunCommand.
	RunCommand(shellCommand string) error
}

// LoopingExecutor supports special casing of looping commands.
type LoopingExecutor interface {
	// RunLoopingCommand works like RunCommand, but invokes leader
	// with the current state again after running the command.
	RunLoopingCommand(shellCommand string) error
}
