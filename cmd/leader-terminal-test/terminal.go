package main

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/Nerdmaster/terminal"
)

// Terminal represents a TTY
type Terminal struct {
	fd            int
	out           io.Writer
	originalState *terminal.State
	keyReader     *terminal.KeyReader
}

// NewTerminalTTY returns a terminal connected to /dev/tty.
func NewTerminalTTY() (*Terminal, error) {
	devTTY, err := os.OpenFile("/dev/tty", os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("NewTerminalTTY: %s", err)
	}
	tty := &Terminal{
		fd:  int(devTTY.Fd()),
		out: devTTY,
	}
	tty.keyReader = terminal.NewKeyReader(devTTY)

	return tty, nil
}

// MakeRaw puts this terminal into raw mode.
func (term *Terminal) MakeRaw() error {
	originalState, err := terminal.MakeRaw(term.fd)
	if err != nil {
		return fmt.Errorf("terminal.GetState: %s", err)
	}
	term.originalState = originalState
	return nil
}

// InputFrom sets up this terminal to read its input from the provided io.Reader.
func (term *Terminal) InputFrom(src io.Reader) *Terminal {
	term.keyReader = terminal.NewKeyReader(src)
	return term
}

// Write implements io.Writer by writing bytes to the underlying terminal.
func (term *Terminal) Write(data []byte) (int, error) {
	return term.out.Write(data)
}

// Restore restores the original terminal state
func (term *Terminal) Restore() error {
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
func (term *Terminal) ReadKey() (rune, error) {
	keypress, err := term.keyReader.ReadKeypress()
	ctrlC := rune('\003')
	if err != nil {
		return ctrlC, nil
	}
	return keypress.Key, nil
}
