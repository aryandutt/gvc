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

func IsWorkingDirClean() (bool, error) {
	index, err := LoadIndex()
	if err != nil {
		return false, fmt.Errorf("failed to load index: %v", err)
	}
	wdMap, err := ScanWorkingDir()
	if err != nil {
		return false, fmt.Errorf("failed to scan working directory: %v", err)
	}
	for _, entry := range *index {
		currentHash, exists := wdMap[entry.Path]
		// Check if the file in index doesn't exist in working directory
		if !exists {
			return false, nil
		}
		// If the blob hash in working directory differs from index, it's not clean
		if currentHash != entry.BlobHash {
			return false, nil
		}
	}
	return true, nil
}

//  Matches the working directory with the commit and the index with the new HEAD.
// TODO: Make this function a transaction to avoid partial updates
func MatchDirectoryWithCommit(commitHash string) error {
	commit, err := GetCommit(commitHash)
	if err != nil {
		return fmt.Errorf("failed to get commit: %v", err)
	}
	treeFiles, err := GetTreeFiles(commit.Tree)
	if err != nil {
		return fmt.Errorf("failed to get tree files: %v", err)
	}
	wdMap, err := ScanWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to scan working dir: %v", err)
	}
	index, err := LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %v", err)
	}
	for path := range wdMap {
		// Skip untracked files
		if _, exists := index.GetEntry(path); !exists {
			continue
		}
		if _, exists := treeFiles[path]; !exists {
			// If the file in the working directory is not in the tree, delete it
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("failed to remove file: %v", err)
			}
			continue
		}
	}

	newIndex := &Index{}
	// Matching tree to index
	for path, hash := range treeFiles {
		// TODO: check the type of file
		// Adding entries to match the index with the commit
		*newIndex = append(*newIndex, IndexEntry{
			Path:     path,
			BlobHash: hash,
			Type: "100644",
		})
		// Modify the files in working dir to match the commit
		data, err := ReadBlobData(hash)
		if err != nil {
			return fmt.Errorf("failed to read blob data: %v", err)
		}
		err = os.WriteFile(path, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file: %v", err)
		}
	}
	
	return newIndex.SaveIndex()
}
