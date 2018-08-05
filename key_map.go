package main

import (
	"bytes"
	"fmt"
)

type Command interface {
	Execute()
}

type KeyMap struct {
	Name string
	Keys map[rune]interface{}
}

func (km *KeyMap) String() string {
	out := bytes.NewBufferString("")
	for key, value := range km.Keys {
		if child, isKeyMap := value.(*KeyMap); isKeyMap {
			fmt.Fprintf(out, "[%c] %s\n\r", key, child.Name)
			continue
		}
		fmt.Fprintf(out, "[%c] %s\n\r", key, value)
	}
	return out.String()
}

func (km *KeyMap) HandleKey(key rune) (*KeyMap, Command) {
	next, found := km.Keys[key]
	if !found {
		return km, nil
	}
	cmd, isCommand := next.(Command)
	keyMap, _ := next.(*KeyMap)

	if isCommand {
		return nil, cmd
	}

	return keyMap, nil
}
