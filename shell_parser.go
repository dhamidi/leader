package main

// ShellParser provides simple parsing methods to detect the lexical
// context in which leader is invoked as part of an interactive shell.
//
// Using ShellParser, leader can detect whether it should run or not.
type ShellParser struct{}

// NewShellParser creates a new parser instance.
func NewShellParser() *ShellParser {
	return &ShellParser{}
}

// InQuotedString returns true if cursor points to a position in line that is within a quoted string.
func (p *ShellParser) InQuotedString(line string, cursor int) bool {
	textUpToCursor := line[:cursor]
	doubleQuotes := 0
	singleQuotes := 0
	escaped := false
	for _, char := range []rune(textUpToCursor) {
		switch char {
		case '"':
			if doubleQuotes == 0 {
				doubleQuotes++
			} else if !escaped {
				doubleQuotes--
			}
		case '\'':
			if singleQuotes == 0 {
				singleQuotes++
			} else {
				singleQuotes--
			}
		case '\\':
			escaped = !escaped
		}
	}
	return doubleQuotes != 0 || singleQuotes != 0
}
