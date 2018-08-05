package main

import "fmt"

type Command interface {
	Execute()
}
type CommandFn func(*MenuState) (Command, error)

type KeyMap struct {
	Name string
	Keys map[rune]interface{}
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

func (km *KeyMap) Merge(other *KeyMap) error {
	for key, entry := range other.Keys {
		if km.Keys[key] == nil {
			km.Keys[key] = entry
			continue
		}
		if entry == nil {
			continue
		}

		entryAsKeyMap, entryIsKeyMap := entry.(*KeyMap)
		thisEntryAsKeyMap, thisEntryIsKeyMap := km.Keys[key].(*KeyMap)
		if entryIsKeyMap && thisEntryIsKeyMap {
			thisEntryAsKeyMap.Merge(entryAsKeyMap)
			continue
		}
		if thisEntryIsKeyMap && !entryIsKeyMap || !thisEntryIsKeyMap && entryIsKeyMap {
			return fmt.Errorf("cannot merge %#v into %#v\n", entry, km.Keys[key])
		}

		km.Keys[key] = entry
	}

	return nil
}
