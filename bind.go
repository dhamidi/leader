package main

import "fmt"

// Bind adds a new key binding.
type Bind struct {
	*Context
	key          string
	boundCommand string
}

// NewBind creates a new bind command to bind key to boundCommand.
func NewBind(context *Context, key string, boundCommand string) *Bind {
	return &Bind{
		Context:      context,
		key:          key,
		boundCommand: boundCommand,
	}
}

// Execute implements Command by adding the key binding encoded by
// this command to the leader configuration file in the current
// directory.
func (cmd *Bind) Execute() error {
	configFile, err := cmd.Files.Open(".leaderrc")
	if err != nil {
		return fmt.Errorf("Bind: open config file: %s", err)
	}
	defer closeIfPossible(configFile)
	config := NewConfig()
	config.Root.Name = nil
	if err := config.ParseJSON(configFile); err != nil {
		return fmt.Errorf("Bind: parse config file: %s", err)
	}
	closeIfPossible(configFile)
	currentConfigMap := config.Root
	for i := 0; i < len(cmd.key)-1; i++ {
		currentConfigMap = currentConfigMap.FindOrAdd(cmd.key[i : i+1])
	}
	currentConfigMap.Keys[cmd.key[0:1]] = &ConfigBinding{
		ShellCommand: &cmd.boundCommand,
	}

	writeToConfigFile, err := cmd.Files.Create(".leaderrc~")
	if err != nil {
		return fmt.Errorf("Bind: failed to create new config file: %s", err)
	}
	if err := config.EncodeJSON(writeToConfigFile); err != nil {
		return fmt.Errorf("Bind: failed to write new config file: %s", err)
	}
	closeIfPossible(writeToConfigFile)
	if err := cmd.Files.Rename(".leaderrc~", ".leaderrc"); err != nil {
		return fmt.Errorf("Bind: failed to install new config file: %s", err)
	}

	return nil
}
