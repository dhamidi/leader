package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestBind_Execute_adds_a_binding_to_local_configuration_by_default(t *testing.T) {
	childRC := `{"keys": {"h": "home", "c": "child"}}`
	context := newTestContextForConfig(t)
	defineTestFile(context, ".leaderrc", childRC)

	main.NewBind(context, "c", "bind").Execute()
	expectedConfig, _ := json.MarshalIndent(map[string]interface{}{
		"keys": map[string]interface{}{
			"h": "home",
			"c": "bind",
		},
	}, "", "  ")
	configFile, err := context.Files.Open(".leaderrc")
	assert.NoError(t, err)
	actualConfig, _ := ioutil.ReadAll(configFile)
	assert.Equal(t, string(expectedConfig), string(actualConfig))

}

func TestBind_Execute_supports_nested_bindings(t *testing.T) {
	childRC := `{"keys": {"h": "home", "c": "child"}}`
	context := newTestContextForConfig(t)
	defineTestFile(context, ".leaderrc", childRC)

	main.NewBind(context, "cc", "bind").Execute()
	expectedConfig, _ := json.MarshalIndent(map[string]interface{}{
		"keys": map[string]interface{}{
			"h": "home",
			"c": map[string]interface{}{
				"keys": map[string]interface{}{
					"c": "bind",
				},
			},
		},
	}, "", "  ")
	configFile, err := context.Files.Open(".leaderrc")
	assert.NoError(t, err)
	actualConfig, _ := ioutil.ReadAll(configFile)
	assert.Equal(t, string(expectedConfig), string(actualConfig))

}

func TestBind_Execute_writes_to_home_leaderrc_if_option_global_is_given(t *testing.T) {
	homeRCPath := os.ExpandEnv("${HOME}/.leaderrc")
	homeRC := `{"keys": {"h": "home", "c": "child"}}`
	context := newTestContextForConfig(t)
	defineTestFile(context, homeRCPath, homeRC)

	main.NewBind(context, "cc", "bind").
		SetGlobal(homeRCPath).
		Execute()
	expectedConfig, _ := json.MarshalIndent(map[string]interface{}{
		"keys": map[string]interface{}{
			"h": "home",
			"c": map[string]interface{}{
				"keys": map[string]interface{}{
					"c": "bind",
				},
			},
		},
	}, "", "  ")
	configFile, err := context.Files.Open(homeRCPath)
	assert.NoError(t, err)
	actualConfig, _ := ioutil.ReadAll(configFile)
	assert.Equal(t, string(expectedConfig), string(actualConfig))

}
