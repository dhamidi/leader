package main_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Execute_loads_less_specific_config_files_first(t *testing.T) {
	homeRC := `{"keys": {"h": "home", "p": "home", "c": "home"}}`
	parentRC := `{"keys": {"h": "home", "p": "parent", "c": "parent"}}`
	childRC := `{"keys": {"h": "home", "p": "parent", "c": "child"}}`
	context := newTestContextForConfig(t)
	defineTestFile(context, os.ExpandEnv("${HOME}/.leaderrc"), homeRC)
	defineTestFile(context, filepath.Clean(os.ExpandEnv("${PWD}/../.leaderrc")), parentRC)
	defineTestFile(context, os.ExpandEnv("${PWD}/.leaderrc"), childRC)

	main.NewLoadConfig(context, os.ExpandEnv("${PWD}"), os.ExpandEnv("${HOME}")).Execute()

	assert.Equal(t, "home", context.CurrentKeyMap.LookupKey('h').Description())
	assert.Equal(t, "parent", context.CurrentKeyMap.LookupKey('p').Description())
	assert.Equal(t, "child", context.CurrentKeyMap.LookupKey('c').Description())
}

func TestLoadConfig_Execute_tries_to_load_home_leaderrc_even_when_outside_of_home_directory(t *testing.T) {
	homeRC := `{"keys": {"h": "home", "c": "home"}}`
	currentRC := `{"keys": {"c": "child"}}`
	context := newTestContextForConfig(t)
	defineTestFile(context, os.ExpandEnv("${HOME}/.leaderrc"), homeRC)
	defineTestFile(context, "/tmp/.leaderrc", currentRC)

	main.NewLoadConfig(context, "/tmp", os.ExpandEnv("${HOME}")).Execute()

	assert.Equal(t, "home", context.CurrentKeyMap.LookupKey('h').Description())
	assert.Equal(t, "child", context.CurrentKeyMap.LookupKey('c').Description())
}
