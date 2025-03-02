package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func Status() error {
	index, err := LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %v", err)
	}

	// headTree
	
	_, err = GetHeadTree()
	if err != nil {
		return fmt.Errorf("failed to get head tree: %v", err)
	}

	wdMap, err := ScanWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to scan working directory: %v", err)
	}


	var _, modifiedFiles, deletedFiles, untrackedFiles []string

	for _, entry := range *index {
		currentHash, exists := wdMap[entry.Path]
		if !exists {
			deletedFiles = append(deletedFiles, entry.Path)
			continue
		}
		// If the blob hash in working directory differs from index, mark as modified
		if currentHash != entry.BlobHash {
			modifiedFiles = append(modifiedFiles, entry.Path)
		}
	}

	for path := range wdMap {
		_, exists := index.GetEntry(path)
		if !exists {
			untrackedFiles = append(untrackedFiles, path)
		}
	}

	fmt.Println("\nChanges not staged for commit:")

	red := color.New(color.FgHiRed).SprintFunc()

	for _, path := range modifiedFiles {
		fmt.Printf("\t%s: %s\n", red("modified"), red(path))
	}
	for _, path := range deletedFiles {
		fmt.Printf("\t%s: %s\n", red("deleted"), red(path))
	}

	fmt.Println("\nUntracked files:")
	for _, path := range untrackedFiles {
		fmt.Printf("\t%s: %s\n", red("untracked"), red(path))
	}

	return nil
}

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
