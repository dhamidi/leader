package main_test

import (
	"bytes"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

func TestConfig_ParseJSON(t *testing.T) {
	exampleConfig := `
{
  "keys": {
    "d": "date",
    "r": {
      "name": "rake",
      "keys": {
        "t": "bundle exec rake test"
      }
    }
  }
}
`
	config := main.NewConfig()
	assert.NoError(t, config.ParseJSON(bytes.NewBufferString(exampleConfig)))
	assert.Equal(t, "root", *config.Root.Name)
	assert.Equal(t, "date", *(config.Root.Keys["d"].ShellCommand))
	assert.Equal(t, "bundle exec rake test",
		*(config.Root.Keys["r"].Child.Keys["t"].ShellCommand),
	)
	assert.Equal(t, "rake", *(config.Root.Keys["r"].Child.Name))
}
