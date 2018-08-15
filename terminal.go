package main

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/Nerdmaster/terminal"
)

// Terminal is a terminal device.
type Terminal interface {
	io.Writer
	MakeRaw() error
	Restore() error
	ReadKey() (rune, error)
}

// TTY represents a TTY and implements Terminal
type TTY struct {
	fd            int
	file          *os.File
	out           io.Writer
	originalState *terminal.State
	keyReader     *terminal.KeyReader
}

// NewTTY returns a terminal connected to /dev/tty.
func NewTTY() (*TTY, error) {
	devTTY, err := os.OpenFile("/dev/tty", os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("NewTTY: %s", err)
	}
	tty := &TTY{
		fd:   int(devTTY.Fd()),
		out:  devTTY,
		file: devTTY,
	}
	tty.keyReader = terminal.NewKeyReader(devTTY)

	return tty, nil
}

// File returns the file object connected to this terminal (or nil if
// this terminal is not connected to a file)
func (term *TTY) File() *os.File {
	_, connectedToFile := term.out.(*os.File)
	if !connectedToFile {
		return nil
	}
	return term.file
}

// MakeRaw puts this terminal into raw mode.
func (term *TTY) MakeRaw() error {
	originalState, err := terminal.MakeRaw(term.fd)
	if err != nil {
		return fmt.Errorf("terminal.GetState: %s", err)
	}
	term.originalState = originalState
	return nil
}

// OutputTo to sets up this terminal to write its output into the provided io.Writer.
func (term *TTY) OutputTo(out io.Writer) *TTY {
	term.out = out
	return term
}

// InputFrom sets up this terminal to read its input from the provided io.Reader.
func (term *TTY) InputFrom(src io.Reader) *TTY {
	term.keyReader = terminal.NewKeyReader(src)
	return term
}

// Write implements io.Writer by writing bytes to the underlying terminal.
func (term *TTY) Write(data []byte) (int, error) {
	return term.out.Write(data)
}

// Restore restores the original terminal state
func (term *TTY) Restore() error {
	if term.originalState == nil {
		return nil
	}
	err := terminal.Restore(term.fd, term.originalState)
	if errno, ok := err.(syscall.Errno); ok && errno == 0 {
		return nil
	}
	return err
}

// ReadKey reads a single key code from the terminal
func (term *TTY) ReadKey() (rune, error) {
	keypress, err := term.keyReader.ReadKeypress()
	ctrlC := rune('\003')
	if err != nil {
		return ctrlC, nil
	}
	return keypress.Key, nil
}
