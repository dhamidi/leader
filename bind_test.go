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
	homeRC := `{"keys": {"h": "home", "c": "home"}}`
	childRC := `{"keys": {"h": "home", "c": "child"}}`
	context := newTestContextForConfig(t)
	defineTestFile(context, os.ExpandEnv("${HOME}/.leaderrc"), homeRC)
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
