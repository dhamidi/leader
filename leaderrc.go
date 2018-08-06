package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ReloadConfig struct {
	State     *MenuState
	Home      string
	StartFrom string
	Verbose   bool
}

func (cmd *ReloadConfig) String() string { return "reload config" }

func (cmd *ReloadConfig) Execute() {
	directory := cmd.StartFrom
	filesToLoad := []string{}
	for {
		filename := filepath.Join(directory, ".leaderrc")
		directory = filepath.Dir(directory)
		filesToLoad = append(filesToLoad, filename)
		if len(directory) < len(cmd.Home) {
			break
		}
	}

	for i := len(filesToLoad); i > 0; i-- {
		filename := filesToLoad[i-1]
		err := LoadLeaderRC(filename, cmd.State)
		if err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(cmd.State.Err, "Error loading rc file: %s\n", err)
			continue
		}
		if cmd.Verbose && err == nil {
			fmt.Printf("Loaded %s\n", filename)
		}
	}
}

type LeaderRCJSON struct {
	Bindings map[string]interface{}
}

func LoadLeaderRCString(rc string, state *MenuState) error {
	src := bytes.NewBufferString(rc)
	rawData := &LeaderRCJSON{}
	if err := json.NewDecoder(src).Decode(rawData); err != nil {
		return err
	}

	return parseKeyBindings(rawData, state)
}

func LoadLeaderRC(filename string, state *MenuState) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	rawData := &LeaderRCJSON{}
	if err := json.NewDecoder(file).Decode(rawData); err != nil {
		return err
	}

	return parseKeyBindings(rawData, state)
}

func parseKeyBindings(rc *LeaderRCJSON, state *MenuState) error {
	result := &KeyMap{Name: "global", Keys: map[rune]interface{}{}}

	for key, keyMapOrCommand := range rc.Bindings {
		keyRune, entry, err := parseKeyBinding(key, keyMapOrCommand, "", state)
		if err != nil {
			return err
		}
		result.Keys[keyRune] = entry
	}

	if state.Root == nil {
		state.Root = result
		return nil
	}

	return state.Root.Merge(result)
}

func parseKeyMap(keyMap map[string]interface{}, path string, state *MenuState) (*KeyMap, error) {
	name := ""
	if keyMap["name"] != nil {
		keyMapName, isString := keyMap["name"].(string)
		if !isString {
			return nil, fmt.Errorf("%s.name needs to be a string", path)
		}
		name = keyMapName
	}
	result := &KeyMap{Name: name, Keys: map[rune]interface{}{}}
	keys, isObject := keyMap["keys"].(map[string]interface{})
	if !isObject {
		return nil, fmt.Errorf("%s needs to be an object", path)
	}
	path = fmt.Sprintf("%s.keys", path)
	for key, keyMapOrCommand := range keys {
		keyRune, entry, err := parseKeyBinding(key, keyMapOrCommand, path, state)
		if err != nil {
			return nil, err
		}
		result.Keys[keyRune] = entry
	}

	return result, nil
}
func parseKeyBinding(key string, keyMapOrCommand interface{}, path string, state *MenuState) (rune, interface{}, error) {
	asCommand, isCommand := keyMapOrCommand.([]interface{})
	asKeyMap, isKeyMap := keyMapOrCommand.(map[string]interface{})
	if !isKeyMap && !isCommand {
		return ' ', nil, fmt.Errorf("%s.%s needs to be an Object or Array", path, key)
	}
	keyRune := keyFromString(key)

	if isCommand {
		commandFn, err := parseCommand(asCommand)
		if err != nil {
			return keyRune, nil, fmt.Errorf("%s.%s%s", path, key, err)
		}
		command, err := commandFn(state)
		return keyRune, command, err
	}

	if isKeyMap {
		keyMap, err := parseKeyMap(asKeyMap, path, state)
		if err != nil {
			return keyRune, nil, fmt.Errorf("%s.%s%s", path, key, err)
		}
		return keyRune, keyMap, nil
	}

	return ' ', nil, fmt.Errorf("%s.%s: %#v not a valid binding", path, key, keyMapOrCommand)
}

func parseCommand(parts []interface{}) (CommandFn, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("Empty command description")
	}
	shellCommandPath := ""
	shellCommandArgs := []string{}
	for i, part := range parts {
		word := part.(string)
		if isBuiltinCommand(word) {
			return parseBuiltinCommand(parts)
		}
		if i == 0 {
			shellCommandPath = word
		} else {
			shellCommandArgs = append(shellCommandArgs, word)
		}
	}

	return func(state *MenuState) (Command, error) {
		return NewShellCommand(
			shellCommandPath,
			shellCommandArgs...,
		).
			InputFrom(state.In).
			RedirectTo(state.Out, state.Err), nil

	}, nil
}

func isBuiltinCommand(word string) bool {
	return word[0] == '<' && word[len(word)-1] == '>'
}

func parseBuiltinCommand(parts []interface{}) (CommandFn, error) {
	asStrings := []string{}
	for _, part := range parts {
		partAsString, ok := part.(string)
		if !ok {
			return nil, fmt.Errorf("%v is not a string", part)
		}
		asStrings = append(asStrings, partAsString)
	}

	commandID := asStrings[0][1 : len(asStrings[0])-1]
	return func(state *MenuState) (Command, error) {
		commandFn := state.BuiltinCommands[commandID]
		if commandFn == nil {
			return nil, fmt.Errorf("unknown builtin command: %s", commandID)
		}
		return commandFn(state)
	}, nil
}
