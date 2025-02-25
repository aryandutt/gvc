package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// AddToStage adds a file or dir to the staging area.
func AddToStage(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if (err != nil) {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	if fileInfo.IsDir() {
		entries, err := os.ReadDir(filePath)
		if err != nil {
			return fmt.Errorf("failed to read directory: %v", err)
		}

		for _, entry := range entries {
			if entry.IsDir() && entry.Name()[0] == '.' {
				continue // Ignore hidden directories
			}
			entryPath := filepath.Join(filePath, entry.Name())
			if err := AddToStage(entryPath); err != nil {
				return err
			}
		}
		return nil
	}

	// 1. Create blob from file content
	blobHash, err := createBlobFromFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to create blob: %v", err)
	}

	// 2. Update index
	index, err := LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %v", err)
	}

	// Remove existing entry if file is already staged
	newIndex := removeEntry(index, filePath)

	fileType := "normal"
	if fileInfo.Mode().Perm()&0111 != 0 {
		fileType = "executable"
	}

	// Add new entry
	*newIndex = append(*newIndex, IndexEntry{
		Path:     filePath,
		BlobHash: blobHash,
		Type:     fileType,
	})

	if err := newIndex.SaveIndex(); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	return nil
}

// Helper: Remove existing entries for a file path
func removeEntry(index *Index, path string) *Index {
    cleanPath := filepath.Clean(path)
    newIndex := &Index{} // Creates a pointer to a new Index
    
    for _, entry := range *index {
        if entry.Path != cleanPath {
			*newIndex = append(*newIndex, entry)
        }
    }
    
    return newIndex
}

// Helper: Create blob from file content
func createBlobFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return CreateObject("blob", content)
}
