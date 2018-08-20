package main_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

// testFileSystem returns files from predefined buffers holding the file's contents
type testFileSystem struct {
	files map[string]string
}

// NewTestFileSystem returns a new, empty test file system.
func NewTestFileSystem() *testFileSystem {
	return &testFileSystem{
		files: map[string]string{},
	}
}

// Open implements main.FileOpener
func (fs *testFileSystem) Open(name string) (io.Reader, error) {
	contents, found := fs.files[name]
	if !found {
		return nil, os.ErrNotExist
	}

	return bytes.NewBufferString(contents), nil
}

// Define sets the contents for a given file in this test file system
func (fs *testFileSystem) Define(path string, contents string) *testFileSystem {
	fs.files[path] = contents
	return fs
}

func TestLoadConfigFile_Execute_merges_key_bindings_from_config_file(t *testing.T) {
	configFile := `
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
`
	keymap := main.NewKeyMap("root")
	context := newTestContext(t, keymap, bytes.NewBufferString(""), nil)
	context.Files.(*testFileSystem).Define(".leaderrc", configFile)
	loadConfig := main.NewLoadConfigFile(context, ".leaderrc")
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
