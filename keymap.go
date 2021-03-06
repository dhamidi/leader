package main

import (
	"fmt"
	"reflect"
	"sort"
)

// Command represents an executable action to modify state.  The
// action can potentially result in an error.
type Command func() error

// DoNothing is a command that never fails and does not change any state.
func DoNothing() error { return nil }

// FailWhenExecuted returns a command that returns the given error when executed.
func FailWhenExecuted(err error) Command {
	return func() error { return err }
}

// KeyMap is a named binding of keys to functions or other keymaps.
type KeyMap struct {
	name     string
	bindings map[rune]*KeyBinding
}

// KeyBinding represents a binding of a key to command or keymap.
type KeyBinding struct {
	key         rune
	command     Command
	loops       bool
	children    *KeyMap
	description string
}

var (
	// UnboundKey is a null-object representing non-existing key bindings.
	//
	// Executing its associated command returns an error.
	UnboundKey = NewKeyBinding('?').Do(DoNothing)
)

// NewKeyBinding returns a new key binding for the given key.  It is
// bound to DoNothing and has no children.
func NewKeyBinding(key rune) *KeyBinding {
	return &KeyBinding{
		key:      key,
		command:  DoNothing,
		children: NewKeyMap(""),
	}
}

// String returns a human readable representation of this key binding.  It does not descend into any child key maps.
func (b *KeyBinding) String() string {
	if b.HasChildren() {
		return fmt.Sprintf("[%c] <keymap %s>", b.key, b.children.Name())
	}

	return fmt.Sprintf("[%c] %s", b.key, b.description)
}

// SetLooping marks this key binding as a looping key binding.  This
// key can be pressed repeatedly to execute the same command
// repeatedly.
func (b *KeyBinding) SetLooping(isLooping bool) *KeyBinding {
	b.loops = isLooping
	return b
}

// IsLooping returns true if this is a looping key binding.
func (b *KeyBinding) IsLooping() bool { return b.loops }

// Describe sets the description for this key binding and returns this binding.
func (b *KeyBinding) Describe(description string) *KeyBinding {
	b.description = description
	return b
}

// HasChildren returns true if any child bindings have been defined for this binding.
func (b *KeyBinding) HasChildren() bool {
	return len(b.children.bindings) > 0
}

// Children returns the keymap associated with this binding.
func (b *KeyBinding) Children() *KeyMap { return b.children }

// Execute runs the command associated with this binding.
func (b *KeyBinding) Execute() error {
	return b.command()
}

// Do associated a command with this key binding and returns this key binding.
func (b *KeyBinding) Do(cmd Command) *KeyBinding {
	b.command = cmd
	return b
}

// IsBoundToCommand returns true if this key binding is not bound to DoNothing.
func (b *KeyBinding) IsBoundToCommand() bool {
	return reflect.ValueOf(b.command).Pointer() != reflect.ValueOf(DoNothing).Pointer()
}

// Key returns the key this key binding is bound to.
func (b *KeyBinding) Key() rune {
	return b.key
}

// Description returns the description of this key binding.
func (b *KeyBinding) Description() string {
	return b.description
}

// NewKeyMap initializes a new keymap named name.
//
// The map has no initial bindings.
//
// Use DefineKey to add bindings to this map.
func NewKeyMap(name string) *KeyMap {
	return &KeyMap{
		name:     name,
		bindings: map[rune]*KeyBinding{},
	}
}

// Name returns the name of this key map
func (m *KeyMap) Name() string {
	return m.name
}

// Rename changes the name of this keymap to name and returns this keymap.
func (m *KeyMap) Rename(newName string) *KeyMap {
	m.name = newName
	return m
}

// Bind adds an empty keybinding this map at key and returns it.
func (m *KeyMap) Bind(key rune) *KeyBinding {
	binding := NewKeyBinding(key)
	m.bindings[key] = binding
	return binding
}

// Set adds the given binding to this key map, merging any child bindings if necessary.
func (m *KeyMap) Set(b *KeyBinding) *KeyMap {
	if m.bindings[b.key] == nil {
		m.bindings[b.key] = b
		return m
	}

	if !b.HasChildren() {
		m.bindings[b.key] = b
		return m
	}

	existingBinding := m.bindings[b.key]
	for _, child := range b.Children().Bindings() {
		existingBinding.Children().Set(child)
	}

	return m
}

// DefineKey binds key to the given command.
//
// It returns the keymap on which this method is called so that method
// call can be chained on this object.
func (m *KeyMap) DefineKey(key rune, cmd Command) *KeyMap {
	m.bindings[key] = NewKeyBinding(key).Do(cmd)
	return m
}

// LookupKey search for a binding for key in this map and returns it.
//
// If no such binding exists, UnboundKey is returned.
func (m *KeyMap) LookupKey(key rune) *KeyBinding {
	binding, found := m.bindings[key]
	if !found {
		return UnboundKey
	}

	return binding
}

// Bindings returns all key bindings in alphabetically ascending order.
func (m *KeyMap) Bindings() []*KeyBinding {
	bindings := []*KeyBinding{}
	for _, binding := range m.bindings {
		bindings = append(bindings, binding)
	}
	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].key < bindings[j].key
	})
	return bindings
}
