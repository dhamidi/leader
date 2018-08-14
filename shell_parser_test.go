package main_test

import (
	"testing"

	"github.com/dhamidi/leader"
	"github.com/stretchr/testify/assert"
)

const (
	printfDate         = `printf "%s\n" "$(date)"`
	escapedDoubleQuote = `printf "\"hello\"\n"`
	printfDateSingle   = `printf '%s\n' "$(date)"`
	bindMixedQuotes    = `bind -x '"\\":leader'`
)

func TestShellParser_InQuotedString_returns_true_if_cursor_is_inside_double_quotes(t *testing.T) {
	parser := main.NewShellParser()
	assert.True(t, parser.InQuotedString(printfDate, len(`printf "`)))
}

func TestShellParser_InQuotedString_returns_false_if_cursor_is_outside_of_double_quotes(t *testing.T) {
	parser := main.NewShellParser()
	assert.False(t, parser.InQuotedString(printfDate, len(`printf `)))
}

func TestShellParser_InQuotedString_returns_true_if_string_contains_escaped_double_quotes(t *testing.T) {
	parser := main.NewShellParser()
	assert.True(t, parser.InQuotedString(escapedDoubleQuote, len(`printf "\"h`)))
}

func TestShellParser_InQuotedString_returns_true_if_cursor_is_inside_single_quotes(t *testing.T) {
	parser := main.NewShellParser()
	assert.True(t, parser.InQuotedString(printfDateSingle, len(`printf '`)))
}

func TestShellParser_InQuotedString_returns_false_if_cursor_is_outside_of_single_quotes(t *testing.T) {
	parser := main.NewShellParser()
	assert.False(t, parser.InQuotedString(printfDateSingle, len(`printf `)))
}

func TestShellParser_InQuotedString_returns_true_if_cursor_is_inside_nested_quotes(t *testing.T) {
	parser := main.NewShellParser()
	assert.True(t, parser.InQuotedString(bindMixedQuotes, len(`bind -x '"\`)))
}
