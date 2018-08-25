package main_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

// writeableTestFile wraps *bytes.Buffer by updating a given string after each call to Write.
type writeableTestFile struct {
	*bytes.Buffer
	updateAfterWrite *string
}

func (w *writeableTestFile) Write(p []byte) (n int, err error) {
	n, err = w.Buffer.Write(p)
	*w.updateAfterWrite = w.Buffer.String()
	return
}

// testFileSystem returns files from predefined buffers holding the file's contents
type testFileSystem struct {
	files []*struct{ name, contents string }
}

// newTestFileSystem returns a new, empty test file system.
func newTestFileSystem() *testFileSystem {
	return &testFileSystem{
		files: []*struct{ name, contents string }{},
	}
}

// Open implements main.FileSystem by returning a handle to an internal buffer.
func (fs *testFileSystem) Open(name string) (io.Reader, error) {
	for _, f := range fs.files {
		if f.name == name {
			return bytes.NewBufferString(f.contents), nil
		}
	}
	return nil, os.ErrNotExist
}

// Create implements main.FileSystem by returning a handle to an
// internal buffer that can be used for writing.
//
// The internal representation of that buffer is updated after every
// call to Write.  This ensures that a reader of the same file in the
// test file system sees changes made by a writer at the point in time
// the reader is opened.
func (fs *testFileSystem) Create(name string) (io.Writer, error) {
	var file *struct{ name, contents string }
	fs.Define(name, "")
	for _, f := range fs.files {
		if f.name == name {
			file = f
			break
		}
	}
	result := &writeableTestFile{
		Buffer:           bytes.NewBufferString(""),
		updateAfterWrite: &file.contents,
	}

	return result, nil
}

// Rename changes the name of the given file.  If the given source
// does not exist, os.ErrNotExist is returned.
func (fs *testFileSystem) Rename(src, dest string) error {
	var destFile, srcFile *struct{ name, contents string }
	var srcIndex int
	for i, f := range fs.files {
		if f.name == src {
			srcFile = f
			srcIndex = i
			continue
		}
		if f.name == dest {
			destFile = f
			continue
		}
	}

	if srcFile == nil {
		return os.ErrNotExist
	}

	if destFile == nil {
		srcFile.name = dest
		return nil
	}

	destFile.contents = srcFile.contents
	if len(fs.files) == srcIndex+1 {
		fs.files = fs.files[0:srcIndex]
	} else {
		fs.files = append(fs.files[0:srcIndex], fs.files[srcIndex+1:]...)
	}
	return nil
}

// Define sets the contents for a given file in this test file system
func (fs *testFileSystem) Define(path string, contents string) *testFileSystem {
	var file struct{ name, contents string }
	for _, f := range fs.files {
		if f.name == path {
			f.contents = contents
			return fs
		}
	}
	file.name = path
	file.contents = contents
	fs.files = append(fs.files, &file)
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
