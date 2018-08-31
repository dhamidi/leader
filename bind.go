package main

import (
	"fmt"
	"os"
)

// Bind adds a new key binding.
type Bind struct {
	*Context
	key          string
	boundCommand string
	file         string
	unbind       bool
}

// NewBind creates a new bind command to bind key to boundCommand.
func NewBind(context *Context, key string, boundCommand string) *Bind {
	return &Bind{
		Context:      context,
		key:          key,
		boundCommand: boundCommand,
		file:         ".leaderrc",
		unbind:       false,
	}
}

// Unbind tells this instance to remove a binding instead of setting it.
func (cmd *Bind) Unbind() *Bind {
	cmd.unbind = true
	return cmd
}

// SetGlobal configures this command instance to write to the global configuration file at path.
func (cmd *Bind) SetGlobal(path string) *Bind {
	cmd.file = path
	return cmd
}

// Execute implements Command by adding the key binding encoded by
// this command to the leader configuration file in the current
// directory.
func (cmd *Bind) Execute() error {
	config := NewConfig()
	config.Root.Name = nil
	configFile, err := cmd.Files.Open(cmd.file)
	if !os.IsNotExist(err) {
		if err != nil {
			return fmt.Errorf("Bind: open config file: %s", err)
		}
		defer closeIfPossible(configFile)
		if err := config.ParseJSON(configFile); err != nil {
			return fmt.Errorf("Bind: parse config file: %s", err)
		}
		closeIfPossible(configFile)
	}
	currentConfigMap := config.Root
	if cmd.unbind {
		for i := 0; i < len(cmd.key)-1; i++ {
			next := currentConfigMap.Keys[string([]rune{rune(cmd.key[i])})]
			if next.Child == nil && i < len(cmd.key)-1 {
				return nil
			}
			currentConfigMap = next.Child
		}
		delete(currentConfigMap.Keys, cmd.key[len(cmd.key)-1:])
	} else {
		for i := 0; i < len(cmd.key)-1; i++ {
			currentConfigMap = currentConfigMap.FindOrAdd(cmd.key[i : i+1])
		}
		currentConfigMap.Keys[cmd.key[len(cmd.key)-1:]] = &ConfigBinding{
			ShellCommand: &cmd.boundCommand,
		}
	}

	tmpFile := cmd.file + "~"
	writeToConfigFile, err := cmd.Files.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("Bind: failed to create new config file: %s", err)
	}
	if err := config.EncodeJSON(writeToConfigFile); err != nil {
		return fmt.Errorf("Bind: failed to write new config file: %s", err)
	}
	closeIfPossible(writeToConfigFile)
	if err := cmd.Files.Rename(tmpFile, cmd.file); err != nil {
		return fmt.Errorf("Bind: failed to install new config file: %s", err)
	}

	return nil
}
