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
