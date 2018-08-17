package main

import (
	"os"
	"path/filepath"
)

// LoadConfig loads all configuration files that are processed by leader.
type LoadConfig struct {
	*Context
	Start string
	Home  string
}

// NewLoadConfig creates a new instance of LoadConfig which loads
// configuration settings into the provided context.
func NewLoadConfig(ctx *Context, start string, home string) *LoadConfig {
	return &LoadConfig{
		Context: ctx,
		Start:   start,
		Home:    home,
	}
}

// Execute examines the current directory and every parent directory for a file called ".leaderrc" and loads that file as a JSON configuration for leader.
//
// It is not an error if such a file does not exist.
//
// After scanning all directories up to the file system's root
// directory $HOME/.leaderrc is loaded if it has not been loaded
// previously already.
func (cmd *LoadConfig) Execute() error {
	currentPath := cmd.Start
	homeRC := filepath.Join(cmd.Home, ".leaderrc")
	files := []string{homeRC}
	for {
		filename := filepath.Join(currentPath, ".leaderrc")
		currentPath = filepath.Dir(currentPath)
		if len(files) > 0 && filename == files[len(files)-1] {
			break
		}
		if filename == homeRC {
			continue
		}
		files = append(files, filename)
	}

	for _, filename := range files {
		loadConfig := NewLoadConfigFile(cmd.Context, filename)
		if err := loadConfig.Execute(); err != nil && !os.IsNotExist(err) {
			cmd.ErrorLogger.Print(err)
		}
	}
	return nil
}
