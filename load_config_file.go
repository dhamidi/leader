package main

import "io"

// LoadConfigFile loads a configuration file in JSON format and
// applies its contents to the current key map in the application
// context.
type LoadConfigFile struct {
	*Context
	Filename string
}

// NewLoadConfigFile creates a command for loading the configuration stored in filename.
func NewLoadConfigFile(context *Context, filename string) *LoadConfigFile {
	return &LoadConfigFile{
		Context:  context,
		Filename: filename,
	}
}

// Execute implements Command by parsing the file named by filename as
// JSON and applying the configuration settings stored within to the
// current keymap.
func (cmd *LoadConfigFile) Execute() error {
	configFile, err := cmd.Files.Open(cmd.Filename)
	if err != nil {
		return err
	}
	defer func() {
		if closer, ok := configFile.(io.Closer); ok {
			closer.Close()
		}
	}()

	config := NewConfig()
	if err := config.ParseJSON(configFile); err != nil {
		return err
	}

	config.MergeIntoKeyMap(cmd.Context, cmd.CurrentKeyMap)

	return nil
}
