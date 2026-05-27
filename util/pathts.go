package util

import (
	"fmt"
	"os"
	"path/filepath"
)

// PathRelativeToCwd returns the relative path of the target from the current working directory.
func PathRelativeToCwd(targetPath string) (string, error) {
	// 1. Get the absolute path of the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// 2. Ensure target path is absolute (handles cases where targetPath is relative)
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute target path: %w", err)
	}

	// 3. Compute the relative path from CWD to the target
	relPath, err := filepath.Rel(cwd, absTarget)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path: %w", err)
	}

	return relPath, nil
}
