package main

import "io"

// closeIfPossible calls Close on the provided object if it implements io.Closer.
func closeIfPossible(closable interface{}) {
	closer, ok := closable.(io.Closer)
	if !ok {
		return
	}
	closer.Close()
}
