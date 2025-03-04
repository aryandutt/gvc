package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ScanWorkingDir scans the working directory and returns a map of file paths to their blob hash.
func ScanWorkingDir() (map[string]string, error) {
	m := make(map[string]string)
	// Walk through the working directory
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		// ignore all directories which start with .
		if d.IsDir() && len(d.Name()) > 1 && d.Name()[0] == '.' {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			// Read the file content
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %v", err)
			}
			hashStr, _ := GetObjectData("blob", data)
			m[path] = hashStr
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk through working directory: %v", err)
	}
	return m, nil
}
