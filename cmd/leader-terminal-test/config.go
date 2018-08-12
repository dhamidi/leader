package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type Config struct {
	Root *ConfigMap
}

type ConfigMap struct {
	Name     *string
	Bindings ConfigBindings
}

// KeyMapName returns a name suitable for use as the name of a key map.
func (c *ConfigMap) KeyMapName() string {
	if c.Name == nil {
		return ""
	}

	return *c.Name
}

// ConfigBindings represents keybindings stored in a config file.
type ConfigBindings map[string]*ConfigBinding

// ConfigBinding is a single entry in ConfigBindings
type ConfigBinding struct {
	ShellCommand *string
	Child        *ConfigMap
}

func (b *ConfigBinding) UnmarshalJSON(data []byte) error {
	generic := (interface{})(nil)
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}

	switch typed := generic.(type) {
	case string:
		b.ShellCommand = &typed
		return nil
	case map[string]interface{}:
		b.Child = new(ConfigMap)
		return json.Unmarshal(data, b.Child)
	default:
		return fmt.Errorf("Cannot parse %T as valid binding", typed)
	}
}

func NewConfig() *Config {
	return &Config{
		Root: NewConfigMap("root"),
	}
}

func NewConfigMap(name string) *ConfigMap {
	return &ConfigMap{
		Name: &name,
	}
}

func (cfg *Config) ParseJSON(in io.Reader) error {
	return json.NewDecoder(in).Decode(cfg.Root)
}

// MergeIntoKeyMap merges the configuration settings in cfg into the provided keymap.
func (cfg *Config) MergeIntoKeyMap(context *Context, keymap *KeyMap) {
	cfg.mergeConfigMap(context, cfg.Root, keymap)
}

func (cfg *Config) mergeConfigMap(context *Context, configMap *ConfigMap, keymap *KeyMap) {
	for key, binding := range configMap.Bindings {
		keyRune := ([]rune(key))[0]
		keyBinding := NewKeyBinding(keyRune)

		if binding.Child != nil {
			existingBinding := keymap.LookupKey(keyRune)
			if existingBinding == UnboundKey {
				existingBinding = keyBinding
				keymap.Set(existingBinding)
			}
			if existingBinding.Children().Name() == "" &&
				binding.Child.Name != nil {
				existingBinding.Children().Rename(*binding.Child.Name)
			}
			cfg.mergeConfigMap(context, binding.Child, existingBinding.Children())
			continue
		}
		if binding.ShellCommand == nil {
			continue
		}

		keyBinding.
			Do(NewRunShellCommand(context, *binding.ShellCommand).Execute).
			Describe(*binding.ShellCommand)
		keymap.Set(keyBinding)
	}
}
