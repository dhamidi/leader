package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/pkg/term/termios"
)

var keyCodes = map[string][]byte{
	"q":        []byte{0x71},
	"w":        []byte{0x77},
	"e":        []byte{0x65},
	"r":        []byte{0x72},
	"t":        []byte{0x74},
	"y":        []byte{0x79},
	"u":        []byte{0x75},
	"i":        []byte{0x69},
	"o":        []byte{0x6f},
	"p":        []byte{0x70},
	"[":        []byte{0x5b},
	"]":        []byte{0x5d},
	"a":        []byte{0x61},
	"s":        []byte{0x73},
	"d":        []byte{0x64},
	"f":        []byte{0x66},
	"g":        []byte{0x67},
	"h":        []byte{0x68},
	"j":        []byte{0x6a},
	"k":        []byte{0x6b},
	"l":        []byte{0x6c},
	";":        []byte{0x3b},
	"'":        []byte{0x27},
	"\\":       []byte{0x5c},
	"z":        []byte{0x7a},
	"x":        []byte{0x78},
	"c":        []byte{0x63},
	"v":        []byte{0x76},
	"b":        []byte{0x62},
	"n":        []byte{0x6e},
	"m":        []byte{0x6d},
	",":        []byte{0x2c},
	".":        []byte{0x2e},
	"/":        []byte{0x2f},
	"Ctrl+q":   []byte{0x11},
	"Ctrl+w":   []byte{0x17},
	"Ctrl+e":   []byte{0x5},
	"Ctrl+r":   []byte{0x12},
	"Ctrl+t":   []byte{0x14},
	"Ctrl+y":   []byte{0x19},
	"Ctrl+u":   []byte{0x15},
	"Ctrl+i":   []byte{0x9},
	"Ctrl+o":   []byte{0xf},
	"Ctrl+p":   []byte{0x10},
	"Ctrl+[":   []byte{0x1b},
	"Ctrl+]":   []byte{0x1d},
	"Ctrl+a":   []byte{0x1},
	"Ctrl+s":   []byte{0x13},
	"Ctrl+d":   []byte{0x4},
	"Ctrl+f":   []byte{0x6},
	"Ctrl+g":   []byte{0x7},
	"Ctrl+h":   []byte{0x8},
	"Ctrl+j":   []byte{0xa},
	"Ctrl+k":   []byte{0xb},
	"Ctrl+l":   []byte{0xc},
	"Ctrl+;":   []byte{0x3b},
	"Ctrl+'":   []byte{0x27},
	"Ctrl+\\":  []byte{0x1c},
	"Ctrl+z":   []byte{0x1a},
	"Ctrl+x":   []byte{0x18},
	"Ctrl+c":   []byte{0x3},
	"Ctrl+v":   []byte{0x16},
	"Ctrl+b":   []byte{0x2},
	"Ctrl+n":   []byte{0xe},
	"Ctrl+m":   []byte{0xd},
	"Ctrl+,":   []byte{0x2c},
	"Ctrl+.":   []byte{0x2e},
	"Ctrl+/":   []byte{0x1f},
	"Shift+q":  []byte{0x51},
	"Shift+w":  []byte{0x57},
	"Shift+e":  []byte{0x45},
	"Shift+r":  []byte{0x52},
	"Shift+t":  []byte{0x54},
	"Shift+y":  []byte{0x59},
	"Shift+u":  []byte{0x55},
	"Shift+i":  []byte{0x49},
	"Shift+o":  []byte{0x4f},
	"Shift+p":  []byte{0x50},
	"Shift+[":  []byte{0x7b},
	"Shift+]":  []byte{0x7d},
	"Shift+a":  []byte{0x41},
	"Shift+s":  []byte{0x53},
	"Shift+d":  []byte{0x44},
	"Shift+f":  []byte{0x46},
	"Shift+g":  []byte{0x47},
	"Shift+h":  []byte{0x48},
	"Shift+j":  []byte{0x4a},
	"Shift+k":  []byte{0x4b},
	"Shift+l":  []byte{0x4c},
	"Shift+;":  []byte{0x3a},
	"Shift+'":  []byte{0x22},
	"Shift+\\": []byte{0x7c},
	"Shift+z":  []byte{0x5a},
	"Shift+x":  []byte{0x58},
	"Shift+c":  []byte{0x43},
	"Shift+v":  []byte{0x56},
	"Shift+b":  []byte{0x42},
	"Shift+n":  []byte{0x4e},
	"Shift+m":  []byte{0x4d},
	"Shift+,":  []byte{0x3c},
	"Shift+.":  []byte{0x3e},
	"Shift+/":  []byte{0x3f},
	"Alt+q":    []byte{0x1b, 0x71},
	"Alt+w":    []byte{0x1b, 0x77},
	"Alt+e":    []byte{0x1b, 0x65},
	"Alt+r":    []byte{0x1b, 0x72},
	"Alt+t":    []byte{0x1b, 0x74},
	"Alt+y":    []byte{0x1b, 0x79},
	"Alt+u":    []byte{0x1b, 0x75},
	"Alt+i":    []byte{0x1b, 0x69},
	"Alt+o":    []byte{0x1b, 0x6f},
	"Alt+p":    []byte{0x1b, 0x70},
	"Alt+[":    []byte{0x1b, 0x5b},
	"Alt+]":    []byte{0x1b, 0x5d},
	"Alt+a":    []byte{0x1b, 0x61},
	"Alt+s":    []byte{0x1b, 0x73},
	"Alt+d":    []byte{0x1b, 0x64},
	"Alt+f":    []byte{0x1b, 0x66},
	"Alt+g":    []byte{0x1b, 0x67},
	"Alt+h":    []byte{0x1b, 0x68},
	"Alt+j":    []byte{0x1b, 0x6a},
	"Alt+k":    []byte{0x1b, 0x6b},
	"Alt+l":    []byte{0x1b, 0x6c},
	"Alt+;":    []byte{0x1b, 0x3b},
	"Alt+'":    []byte{0x1b, 0x27},
	"Alt+\\":   []byte{0x1b, 0x5c},
	"Alt+z":    []byte{0x1b, 0x7a},
	"Alt+x":    []byte{0x1b, 0x78},
	"Alt+c":    []byte{0x1b, 0x63},
	"Alt+v":    []byte{0x1b, 0x76},
	"Alt+b":    []byte{0x1b, 0x62},
	"Alt+n":    []byte{0x1b, 0x6e},
	"Alt+m":    []byte{0x1b, 0x6d},
	"Alt+,":    []byte{0x1b, 0x2c},
	"Alt+.":    []byte{0x1b, 0x2e},
	"Alt+/":    []byte{0x1b, 0x2f},
}

func keyFromBytes(b []byte) string {
key:
	for keysym, byteseq := range keyCodes {
		for i := range byteseq {
			if b[i] != byteseq[i] {
				continue key
			}
		}
		return keysym
	}

	return "unknown"
}
func rawTerminal() func() {
	fd := uintptr(os.Stdin.Fd())
	newState := new(syscall.Termios)
	if err := termios.Tcgetattr(fd, newState); err != nil {
		panic(err)
	}

	oldState := *newState
	termios.Cfmakeraw(newState)
	if err := termios.Tcsetattr(fd, termios.TCSANOW, newState); err != nil {
		panic(err)
	}

	return func() {
		termios.Tcsetattr(fd, termios.TCSANOW, &oldState)
	}

}

func generateKeyDefs() {
	restoreTerminal := rawTerminal()
	keys := `qwertyuiop[]asdfghjkl;'\zxcvbnm,./`
	modifiers := []string{"", "Ctrl", "Shift", "Alt"}
	output, _ := os.Create("keydefs.go")
	for _, modifier := range modifiers {
		for _, key := range []rune(keys) {
			keyString := string([]rune{key})
			if modifier != "" {
				keyString = modifier + "+" + keyString
			}
			fmt.Printf("Press %s\n\r", keyString)

			keyBuf := make([]byte, 8)
			n, err := os.Stdin.Read(keyBuf)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				break
			}
			keySeq := []byte{}
			for j := 0; j < n; j++ {
				keySeq = append(keySeq, keyBuf[j])
			}
			fmt.Fprintf(output, "%q: %#v\n", keyString, keySeq)
		}
	}
	restoreTerminal()
}

func main() {
	restoreTerminal := rawTerminal()
	for {
		keyBuf := make([]byte, 8)
		_, err := os.Stdin.Read(keyBuf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			break
		}
		keySym := keyFromBytes(keyBuf)
		fmt.Printf("You pressed: %s\n\r", keySym)
	}
	restoreTerminal()
}
