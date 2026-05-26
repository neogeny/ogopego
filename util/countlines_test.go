package util

import "testing"

func TestCountLines_Empty(t *testing.T) {
	got := CountLines("")
	want := LineCounts{0, 0, 0, 0}
	if got != want {
		t.Errorf("CountLines('') = %+v, want %+v", got, want)
	}
}

func TestCountLines_TrailingNewline(t *testing.T) {
	got := CountLines("hello\n")
	want := LineCounts{Total: 2, Blank: 1, Comment: 0, Code: 1}
	if got != want {
		t.Errorf("CountLines('hello\\n') = %+v, want %+v", got, want)
	}
}

func TestCountLines_NoTrailingNewline(t *testing.T) {
	got := CountLines("hello")
	want := LineCounts{Total: 1, Blank: 0, Comment: 0, Code: 1}
	if got != want {
		t.Errorf("CountLines('hello') = %+v, want %+v", got, want)
	}
}

func TestCountLines_DefaultCommentPrefix(t *testing.T) {
	s := "code\n// comment\n  // indented comment\nnot a comment\n"
	got := CountLines(s)
	want := LineCounts{Total: 5, Blank: 1, Comment: 2, Code: 2}
	if got != want {
		t.Errorf("CountLines(...) = %+v, want %+v", got, want)
	}
}

func TestCountLines_WhitespaceOnlyLine(t *testing.T) {
	got := CountLines("code\n   \nmore")
	want := LineCounts{Total: 3, Blank: 1, Comment: 0, Code: 2}
	if got != want {
		t.Errorf("CountLines(...) = %+v, want %+v", got, want)
	}
}

func TestCountLines_MultiCharPrefix(t *testing.T) {
	s := "x := 1\n--[[ block ]]\ny := 2\n"
	got := CountLines(s, "--")
	want := LineCounts{Total: 4, Blank: 1, Comment: 1, Code: 2}
	if got != want {
		t.Errorf("CountLines(..., \"--\") = %+v, want %+v", got, want)
	}
}

func TestCountLines_OnlyNewline(t *testing.T) {
	got := CountLines("\n")
	want := LineCounts{Total: 2, Blank: 2, Comment: 0, Code: 0}
	if got != want {
		t.Errorf("CountLines('\\n') = %+v, want %+v", got, want)
	}
}

func TestCountLines_CRLF(t *testing.T) {
	s := "line1\r\nline2\r\n"
	got := CountLines(s)
	want := LineCounts{Total: 3, Blank: 1, Comment: 0, Code: 2}
	if got != want {
		t.Errorf("CountLines('line1\\r\\nline2\\r\\n') = %+v, want %+v", got, want)
	}
}

func TestCountLines_ShebangIsCode(t *testing.T) {
	s := "#!/bin/sh\necho hi\n"
	got := CountLines(s)
	want := LineCounts{Total: 3, Blank: 1, Comment: 0, Code: 2}
	if got != want {
		t.Errorf("CountLines(...) = %+v, want %+v", got, want)
	}
}

func TestCountLines_HashComment(t *testing.T) {
	s := "# shell script\necho hello\n"
	got := CountLines(s, "#")
	want := LineCounts{Total: 3, Blank: 1, Comment: 1, Code: 1}
	if got != want {
		t.Errorf("CountLines(..., \"#\") = %+v, want %+v", got, want)
	}
}
