package main

import (
	"bytes"
	"io"
)

// View is a component that can display data on a terminal.
type View interface {
	Render(out io.Writer) error
}

// RenderViewToString renders view into a string.
func RenderViewToString(view View) (string, error) {
	out := bytes.NewBufferString("")
	err := view.Render(out)
	return out.String(), err
}

// MustRenderViewToString works like RenderViewToString but panics in case of an error.
func MustRenderViewToString(view View) string {
	out, err := RenderViewToString(view)
	if err != nil {
		panic(err)
	}
	return out
}
