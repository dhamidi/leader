package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// ErrorLogger is an error handler that logs errors on a given io.Writer
type ErrorLogger struct {
	logger *log.Logger
}

// NewErrorLogger creates a new error logger that logs errors on the given io.Writer
func NewErrorLogger(out io.Writer) *ErrorLogger {
	return &ErrorLogger{
		logger: log.New(out, "", 0),
	}
}

// Must runs f and calls Fatal if f returned a non-nil error.
func (e *ErrorLogger) Must(f func() error) {
	if err := f(); err != nil {
		e.Fatal(err)
	}
}

// Fatal panics with err
func (e *ErrorLogger) Fatal(err error) { panic(err) }

// Print logs the provided error
func (e *ErrorLogger) Print(err error) {
	e.logger.Printf("error: %s", err)
}

func main() {
	errorHandler := NewErrorLogger(os.Stderr)
	defer func() {
		v := recover()
		if err, ok := v.(error); ok {
			errorHandler.Print(err)
		}
	}()
	tty, err := NewTerminalTTY()
	if err != nil {
		errorHandler.Fatal(err)
	}
	errorHandler.Must(tty.MakeRaw)
	key, err := tty.ReadKey()
	if err != nil {
		errorHandler.Fatal(err)
	}
	errorHandler.Must(tty.Restore)
	fmt.Printf("Key: %c\n", key)
}
