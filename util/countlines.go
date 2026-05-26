package util

import "strings"

// LineCounts holds SLOC (Source Lines of Code) statistics.
type LineCounts struct {
	Total   int
	Blank   int
	Comment int
	Code    int
}

// CountLines counts lines in s with Editor View semantics.
// cmtstr specifies the comment prefix (default: "//").
// A line is blank if it is empty or whitespace-only.
// A line is a comment if its first non-whitespace characters start with cmtstr.
// Editor View: if s ends with '\n' or '\r', an additional blank ghost line
// is counted (the cursor position after the last line terminator).
func CountLines(s string, cmtstrs ...string) LineCounts {
	cmtstr := "//"
	if len(cmtstrs) > 0 {
		cmtstr = cmtstrs[0]
	}

	var totl, blnk, cmnt, code int
	for line := range strings.Lines(s) {
		totl++
		rest := StripLeft(line)
		if rest == "" {
			blnk++
		} else if strings.HasPrefix(rest, cmtstr) {
			cmnt++
		} else {
			code++
		}
	}

	// Editor View: text ending with a line terminator has a ghost blank line.
	if len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		totl++
		blnk++
	}

	return LineCounts{Total: totl, Blank: blnk, Comment: cmnt, Code: code}
}
