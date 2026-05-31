package newlines

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func eq(t *testing.T, label string, got, want int) {
	t.Helper()
	assert.Equal(t, want, got, label)
}

func TestTakeLinebreakLen(t *testing.T) {
	eq(t, "empty line basic", TakeLinebreakLen("  \nrule", 0), 3)
	eq(t, "empty line crlf", TakeLinebreakLen("  \r\nrule", 0), 4)
	eq(t, "content fails", TakeLinebreakLen("content\n", 0), -1)
	eq(t, "eot", TakeLinebreakLen("", 0), 0)
}

func TestBlankLine(t *testing.T) {
	eq(t, "single nl is blank at eot", TakeBlankLineLen("\n", 0), 1)
	eq(t, "empty input is blank", TakeBlankLineLen("", 0), 0)
	eq(t, "standard blank", TakeBlankLineLen("\n\n", 0), 2)
	eq(t, "whitespace blank", TakeBlankLineLen("  \n  \n", 0), 6)
	eq(t, "non-blank content", TakeBlankLineLen("rule\n\n", 0), -1)
	eq(t, "first blank second content", TakeBlankLineLen("\nrule", 0), -1)
	eq(t, "two blanks then rule", TakeBlankLineLen("  \n  \nrule", 0), 6)
	eq(t, "fails on single", TakeBlankLineLen("  \nrule", 0), -1)
	eq(t, "trailing spaces", TakeBlankLineLen("  \n  ", 0), 5)
	eq(t, "two blanks at end", TakeBlankLineLen(" \n ", 0), 3)
	eq(t, "multiple blanks", TakeBlankLineLen("  \n\n\n", 0), 4)
}

func TestDedent(t *testing.T) {
	eq(t, "standard dedent", TakeDedentLen("  \nrule", 0), 3)
	eq(t, "not a dedent", TakeDedentLen(" \n  indented", 0), -1)
	eq(t, "dedent at eot", TakeDedentLen(" \n", 0), 2)
	eq(t, "dedent empty", TakeDedentLen("", 0), 0)
	eq(t, "dedent windows", TakeDedentLen(" \r\ncontent", 0), 3)
}
