// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: MIT OR Apache-2.0

package json

type JsonError struct {
	Message string
}

func (e *JsonError) Error() string {
	return e.Message
}

func NewJsonError(msg string) *JsonError {
	return &JsonError{Message: msg}
}
