package main

import (
	"io"
	"os"
)

// FileSystem defines the necessary operations for opening and renaming files.
//
// Any error that are returned should be of the same set of errors returned by os.OpenFile and os.Rename.
type FileSystem interface {
	Open(name string) (io.Reader, error)
	Create(name string) (io.Writer, error)
	Rename(src, dest string) error
}

// OnDiskFileSystem is a filesystem that gives access to files on
// disk.  It is mainly a wrapper around file related functions from
// the os package.
type OnDiskFileSystem struct{}

// NewOnDiskFileSystem creates a new instance of this filesystem.
func NewOnDiskFileSystem() *OnDiskFileSystem {
	return &OnDiskFileSystem{}
}

// Open implements FileSystem using os.Open
func (fs *OnDiskFileSystem) Open(name string) (io.Reader, error) {
	return os.Open(name)
}

// Create implements FileSystem using os.Create
func (fs *OnDiskFileSystem) Create(name string) (io.Writer, error) {
	return os.Create(name)
}

// Rename implements FileSystem using os.Rename
func (fs *OnDiskFileSystem) Rename(src, dest string) error {
	return os.Rename(src, dest)
}
