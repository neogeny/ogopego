// Copyright (c) 2026 Juancarlo Añez (apalala@gmail.com)
// SPDX-License-Identifier: Apache-2.0

package input

import "github.com/neogeny/ogopego/config"

// Cfg is an alias for the project's configuration type.
type Cfg = config.Cfg

// Configurable is an alias for the configuration interface used by cursors.
type Configurable = config.Configurable

// Sample is a tiny helper used in tests and examples.
func Sample() string {
	return "foo"
}
