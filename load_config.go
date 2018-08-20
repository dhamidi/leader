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

// Execute scans all directories from the current directory up to the
// root directory for files called ".leaderrc".  If such files exist,
// they are parsed as JSON sorted by distance: files further away from
// $HOME/.leaderrc are parsed later than files closer to
// $HOME/.leaderrc.
func (cmd *LoadConfig) Execute() error {
	currentPath := cmd.Start
	homeRC := filepath.Join(cmd.Home, ".leaderrc")
	homeRCAdded := false
	files := []string{}
	for {
		filename := filepath.Join(currentPath, ".leaderrc")
		currentPath = filepath.Dir(currentPath)
		if len(files) > 0 && filename == files[len(files)-1] {
			break
		}
		if filename == homeRC {
			homeRCAdded = true
		}
		files = append(files, filename)
	}

	if !homeRCAdded {
		files = append(files, homeRC)
	}

	for i := len(files) - 1; i >= 0; i-- {
		loadConfig := NewLoadConfigFile(cmd.Context, files[i])
		if err := loadConfig.Execute(); err != nil && !os.IsNotExist(err) {
			cmd.ErrorLogger.Print(err)
		}
	}
	return nil
}
