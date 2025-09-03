package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func ResolvePath(inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	if strings.HasPrefix(inputPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory %w", err)
		}

		if inputPath == "~" {
			inputPath = homeDir
		} else if strings.HasPrefix(inputPath, "~/") || (runtime.GOOS == "windows" && strings.HasPrefix(inputPath, "~\\")) {
			inputPath = filepath.Join(homeDir, inputPath[2:])
		}
	}

	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	cleanPath := filepath.Clean(absPath)

	return cleanPath, nil
}
