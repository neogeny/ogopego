// Package newlines provides platform-agnostic line-boundary detection
// for indentation-based grammar parsing. Ported from TatSu (Python) and TieXiu (Rust).
package newlines

import "strings"

// TakeNonNewlineWhitespaceLen returns the byte length of leading whitespace
// on the current line. It stops at the first newline or non-whitespace character.
func TakeNonNewlineWhitespaceLen(text string, start int) int {
	n := len(text)
	if start >= n {
		return 0
	}
	i := start
	for i < n {
		c := text[i]
		if c == '\n' || (c == '\r' && i+1 < n && text[i+1] == '\n') {
			break
		}
		if c != ' ' && c != '\t' && c != '\r' && c != '\f' && c != '\v' {
			break
		}
		i++
	}
	return i - start
}

// TakeLinebreakLen detects a whitespace-only line or end of text.
// Returns the byte length consumed (including the line terminator),
// or -1 if the line at start contains non-whitespace content.
func TakeLinebreakLen(text string, start int) int {
	n := len(text)
	if start >= n {
		return 0
	}

	nl := strings.IndexByte(text[start:], '\n')
	var endOfLine int
	if nl == -1 {
		endOfLine = n
	} else {
		endOfLine = start + nl
	}

	for i := start; i < endOfLine; i++ {
		c := text[i]
		if c != ' ' && c != '\t' && c != '\r' && c != '\f' && c != '\v' {
			return -1
		}
	}

	if nl != -1 {
		return start + nl - start + 1 // +1 for the \n (skips \r if present via the loop above)
	}
	return n - start
}

// BlankLine matches two consecutive whitespace-only lines.
// Returns the total byte length consumed, or -1 if not found.
func BlankLine(text string, start int) int {
	off1 := TakeLinebreakLen(text, start)
	if off1 < 0 {
		return -1
	}
	off2 := TakeLinebreakLen(text, start+off1)
	if off2 < 0 {
		return -1
	}
	return off1 + off2
}

// IndentLen returns the byte length of leading whitespace on the line at start.
// Returns 0 if start is at or past end of text (EOT is treated as zero margin).
func IndentLen(text string, start int) int {
	n := len(text)
	if start >= n {
		return 0
	}

	nl := strings.IndexByte(text[start:], '\n')
	var searchLimit int
	if nl == -1 {
		searchLimit = n
	} else {
		searchLimit = start + nl
		for searchLimit > start && text[searchLimit-1] == '\r' {
			searchLimit--
		}
	}

	for i := start; i < searchLimit; i++ {
		c := text[i]
		if c != ' ' && c != '\t' && c != '\r' && c != '\f' && c != '\v' {
			return i - start
		}
	}

	return -1
}

// Dedent detects a line break that returns to the zero margin.
// Returns the byte offset consumed, or -1 if not found.
func Dedent(text string, start int) int {
	offset := TakeLinebreakLen(text, start)
	if offset < 0 {
		return -1
	}
	if IndentLen(text, start+offset) == 0 {
		return offset
	}
	return -1
}
