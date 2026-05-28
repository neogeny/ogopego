package util

import (
	"iter"
	"os"

	"golang.org/x/term"
)

// Either returns userVal if it is non-zero, otherwise defaultVal.
func Either[T comparable](userVal, defaultVal T) T {
	var zero T
	if userVal != zero {
		return userVal
	}
	return defaultVal
}

// EitherSlice returns userVal if non-zero, otherwise defaultVal.
func EitherSlice[T any](userVal, defaultVal []T) []T {
	if userVal != nil {
		return userVal
	}
	return defaultVal
}

func TermSize() (int, int) {
	tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err != nil {
		tty = os.Stderr
	} else {
		defer tty.Close()
	}
	cols, rows, err := term.GetSize(int(tty.Fd()))
	if err != nil {
		cols = 88
		rows = 25
	}
	return cols, rows
}

func Chunks[T any](slice []T, chunkSize int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if chunkSize <= 0 {
			return
		}

		for i := 0; i < len(slice); i += chunkSize {
			end := i + chunkSize
			if end > len(slice) {
				end = len(slice)
			}

			// Yield the current batch window. If the loop body
			// breaks early, yield returns false, and we stop.
			if !yield(slice[i:end]) {
				return
			}
		}
	}
}
