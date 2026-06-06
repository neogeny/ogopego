package main

import (
	"fmt"
	"unicode"
)

// scanIdentifier takes the raw input string and current byte offset
func scanIdentifier(input string, start int) (string, int) {
	// 1. Check if the start character is valid using a rune conversion
	// Note: We use the first character at the current offset
	r := []rune(input[start:])[0]
	if !isIdentifierStart(r) {
		return "", start
	}

	// 2. Track the cursor as a byte offset
	cursor := start

	// 3. Loop through the string by byte index
	// We convert the slice to runes only where necessary for Unicode checks
	for cursor < len(input) {
		// Get the current character as a rune
		r := []rune(input[cursor:])[0]

		if !isIdentifierPart(r) {
			break
		}

		// Move cursor forward by the size of the current character
		// Using len(string(r)) ensures we handle UTF-8 multi-byte chars
		cursor += len(string(r))
	}

	// Return the identifier string and the new offset
	return input[start:cursor], cursor
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isIdentifierPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func main() {
	input := "my_variable123"
	token, next := scanIdentifier(input, 0)
	fmt.Printf("Token: %q, Next Offset: %d\n", token, next)
}

