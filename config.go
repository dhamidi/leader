package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Config struct {
	Root *ConfigMap
}

type ConfigMap struct {
	Name        *string        `json:"name,omitempty"`
	LoopingKeys []string       `json:"loopingKeys,omitempty"`
	Keys        ConfigBindings `json:"keys,omitempty"`
}

// FindOrAdd returns the config map associated with key in Keys.
//
// If no ConfigMap exists under the given key, a new anonymous
// ConfigMap is inserted and returned.
//
// If the binding at key does not refer to a ConfigMap, a new
// ConfigMap is constructed which overrides the existing binding.
func (c *ConfigMap) FindOrAdd(key string) *ConfigMap {
	binding, found := c.Keys[key]
	if !found {
		c.Keys[key] = &ConfigBinding{}
		binding = c.Keys[key]
	}

	if binding.Child == nil {
		binding.ShellCommand = nil
		binding.Child = &ConfigMap{
			Keys: ConfigBindings{},
		}
	}

	return binding.Child
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

func (b *ConfigBinding) MarshalJSON() ([]byte, error) {
	if b.Child != nil {
		return json.Marshal(*b.Child)
	}
	return json.Marshal(*b.ShellCommand)
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
		Keys: ConfigBindings{},
	}
}

func (cfg *Config) EncodeJSON(out io.Writer) error {
	marshaled, err := json.MarshalIndent(cfg.Root, "", "  ")
	if err != nil {
		return err
	}
	_, err = io.Copy(out, bytes.NewBuffer(marshaled))
	return err
}

func (cfg *Config) ParseJSON(in io.Reader) error {
	return json.NewDecoder(in).Decode(cfg.Root)
}

// MergeIntoKeyMap merges the configuration settings in cfg into the provided keymap.
func (cfg *Config) MergeIntoKeyMap(context *Context, keymap *KeyMap) {
	cfg.mergeConfigMap(context, cfg.Root, keymap)
}

func (cfg *Config) mergeConfigMap(context *Context, configMap *ConfigMap, keymap *KeyMap) {
	isLoopingKey := func(key rune) bool {
		for _, loopingKey := range configMap.LoopingKeys {
			if key == []rune(loopingKey)[0] {
				return true
			}
		}
		return false
	}
	for key, binding := range configMap.Keys {
		keyRune := ([]rune(key))[0]
		keyBinding := NewKeyBinding(keyRune)
		keyBinding.SetLooping(isLoopingKey(keyRune))
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
