// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package json

// JsonError represents an error that occurs during JSON processing.
type JsonError struct {
	Message string
}

func (e *JsonError) Error() string {
	return e.Message
}

// NewJsonError creates a new JsonError with the given message.
func NewJsonError(msg string) *JsonError {
	return &JsonError{Message: msg}
}
