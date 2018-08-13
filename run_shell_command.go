package main

// RunShellCommand is a command for running shell commands in the current context.
type RunShellCommand struct {
	*Context
	ShellCommand string
}

// NewRunShellCommand creates a command for running shellCommand in ctx.
func NewRunShellCommand(ctx *Context, shellCommand string) *RunShellCommand {
	return &RunShellCommand{
		Context:      ctx,
		ShellCommand: shellCommand,
	}
}

// Execute runs this command.
func (cmd *RunShellCommand) Execute() error {
	if loopingExecutor, canLoop := cmd.Executor.(LoopingExecutor); cmd.CurrentBindingIsLooping() && canLoop {
		return loopingExecutor.RunLoopingCommand(cmd.ShellCommand)
	}
	return cmd.Executor.RunCommand(cmd.ShellCommand)
}
