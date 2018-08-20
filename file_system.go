package main

import (
	"io"
	"os"
)

// FileOpener defines the necessary operations for opening a file for reading on disk.
//
// Any error that are returned should be of the same set of errors returned by os.Open.
type FileOpener interface {
	Open(name string) (io.Reader, error)
}

// OnDiskFileSystem is a filesystem that gives access to files on
// disk.  It is mainly a wrapper around file related functions from
// the os package.
type OnDiskFileSystem struct{}

// NewOnDiskFileSystem creates a new instance of this filesystem.
func NewOnDiskFileSystem() *OnDiskFileSystem {
	return &OnDiskFileSystem{}
}

// Open implements FileOpener using os.Open
func (fs *OnDiskFileSystem) Open(name string) (io.Reader, error) {
	return os.Open(name)
}
