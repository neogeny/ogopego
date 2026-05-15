package newlines

import "testing"

func eq(t *testing.T, label string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %d, want %d", label, got, want)
	}
}

func TestTakeLinebreakLen(t *testing.T) {
	eq(t, "empty line basic", TakeLinebreakLen("  \nrule", 0), 3)
	eq(t, "empty line crlf", TakeLinebreakLen("  \r\nrule", 0), 4)
	eq(t, "content fails", TakeLinebreakLen("content\n", 0), -1)
	eq(t, "eot", TakeLinebreakLen("", 0), 0)
}

func TestBlankLine(t *testing.T) {
	eq(t, "single nl is blank at eot", BlankLine("\n", 0), 1)
	eq(t, "empty input is blank", BlankLine("", 0), 0)
	eq(t, "standard blank", BlankLine("\n\n", 0), 2)
	eq(t, "whitespace blank", BlankLine("  \n  \n", 0), 6)
	eq(t, "non-blank content", BlankLine("rule\n\n", 0), -1)
	eq(t, "first blank second content", BlankLine("\nrule", 0), -1)
	eq(t, "two blanks then rule", BlankLine("  \n  \nrule", 0), 6)
	eq(t, "fails on single", BlankLine("  \nrule", 0), -1)
	eq(t, "trailing spaces", BlankLine("  \n  ", 0), 5)
	eq(t, "two blanks at end", BlankLine(" \n ", 0), 3)
	eq(t, "multiple blanks", BlankLine("  \n\n\n", 0), 4)
}

func TestDedent(t *testing.T) {
	eq(t, "standard dedent", Dedent("  \nrule", 0), 3)
	eq(t, "not a dedent", Dedent(" \n  indented", 0), -1)
	eq(t, "dedent at eot", Dedent(" \n", 0), 2)
	eq(t, "dedent empty", Dedent("", 0), 0)
	eq(t, "dedent windows", Dedent(" \r\ncontent", 0), 3)
}
