package util

import (
	"fmt"
	"runtime/debug"
)

func GetVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	// 1. If installed via 'go install github.com/user/repo@v0.1.0'
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	// 2. If built locally from source (e.g., 'go build'), dig into the VCS settings
	var revision string
	var dirty bool

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			dirty = setting.Value == "true"
		}
	}

	if revision != "" {
		shortRev := revision
		if len(shortRev) > 7 {
			shortRev = shortRev[:7]
		}
		if dirty {
			return fmt.Sprintf("devel-%s-dirty", shortRev)
		}
		return fmt.Sprintf("devel-%s", shortRev)
	}

	return "unknown"
}
