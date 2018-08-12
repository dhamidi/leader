package main_test

import (
	"bytes"
	"testing"

	"github.com/dhamidi/leader/cmd/leader-terminal-test"
	"github.com/stretchr/testify/assert"
)

func TestConfig_ParseJSON(t *testing.T) {
	exampleConfig := `
{
  "bindings": {
    "d": "date",
    "r": {
      "name": "rake",
      "bindings": {
        "t": "bundle exec rake test"
      }
    }
  }
}
`
	config := main.NewConfig()
	assert.NoError(t, config.ParseJSON(bytes.NewBufferString(exampleConfig)))
	assert.Equal(t, "root", *config.Root.Name)
	assert.Equal(t, "date", *(config.Root.Bindings["d"].ShellCommand))
	assert.Equal(t, "bundle exec rake test",
		*(config.Root.Bindings["r"].Child.Bindings["t"].ShellCommand),
	)
	assert.Equal(t, "rake", *(config.Root.Bindings["r"].Child.Name))
}
