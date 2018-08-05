package main

type Command interface {
	Execute()
}
type CommandFn func(*MenuState) Command

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
