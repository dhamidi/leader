package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

// temporaryFileWithContent creates a temporary file with the provided content.
func temporaryFileWithContent(content string) (*os.File, error) {
	tempfile, err := ioutil.TempFile("", "leader-test")
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(tempfile, "%s", content)
	tempfile.Seek(0, 0)
	return tempfile, nil
}

func TestLoadConfigFile_Execute_merges_key_bindings_from_config_file(t *testing.T) {
	configFile, err := temporaryFileWithContent(`
{
  "keys": {
    "d": "date",
    "g": {
      "name": "go",
      "loopingKeys": ["t"],
      "keys": {
        "t": "go test -v ."
      }
    }
  }
}
`)
	assert.NoError(t, err, "creating temporary config file failed")
	defer os.Remove(configFile.Name())

	keymap := main.NewKeyMap("root")
	context := newTestContext(t, keymap, bytes.NewBufferString(""), nil)

	loadConfig := main.NewLoadConfigFile(context, configFile.Name())
	keymap.Bind('d').Describe("do nothing")
	keymap.Bind('g').Children().Bind('t').Describe("go test .")
	assert.NoError(t, loadConfig.Execute(), "loadConfig.Execute()")

	keyD := keymap.LookupKey('d')
	keyG := keymap.LookupKey('g')
	assert.Equal(t, "[d] date", keyD.String())
	assert.Equal(t, "[g] <keymap go>", keyG.String())
	keyGT := keyG.Children().LookupKey('t')
	assert.True(t, keyGT.IsLooping(), "keyGT.IsLooping()")
	assert.Equal(t, "[t] go test -v .", keyGT.String())
}
