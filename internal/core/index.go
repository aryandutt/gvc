package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// IndexEntry represents a file in the staging area.
type IndexEntry struct {
	Path     string `json:"path"`     // File path relative to repo root
	BlobHash string `json:"blobHash"` // SHA-1 hash of the blob
	Type     string `json:"type"`     // Type of file (executable or regular)
}

// Index is the staging area (list of entries).
type Index []IndexEntry

// internal/core/index.go

const indexPath = ".gvc/index.json"

// LoadIndex reads the staging area from .gvc/index.json.
func LoadIndex() (*Index, error) {
	data, err := os.ReadFile(indexPath)
	if os.IsNotExist(err) {
		return &Index{}, nil // Return empty index if file doesn't exist
	} else if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return &Index{}, nil // Return empty index if data is empty
	}

	var index Index
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}
	return &index, nil
}

// SaveIndex writes the staging area to .gvc/index.json.
func (idx *Index) SaveIndex() error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(indexPath, data, 0644)
}

// GetEntry returns the index entry for a given path.
func (idx *Index) GetEntry(path string) (*IndexEntry, bool) {
	for _, entry := range *idx {
		if entry.Path == path {
			return &entry, true
		}
	}
	return nil, false
}

// Returns a list of files that are different between the index and the HEAD commit.
func (index *Index) CompareToHead() ([]string, error) {
    headTreeHash, err := GetHeadTree()
    if err != nil {
        return nil, fmt.Errorf("failed to get head tree: %v", err)
    }

    treeFiles, err := GetTreeFiles(headTreeHash)
    if err != nil {
        return nil, fmt.Errorf("failed to get tree files: %v", err)
    }

    var stagedChanges []string

    for _, entry := range *index {
        headBlob, exists := treeFiles[entry.Path]
        if !exists {
            // File in index is not in HEAD → file added
            stagedChanges = append(stagedChanges, fmt.Sprintf("added: %s", entry.Path))
        } else if headBlob != entry.BlobHash {
            // File exists but blob hash is different → file modified
            stagedChanges = append(stagedChanges, fmt.Sprintf("modified: %s", entry.Path))
        }
    }

    // Optionally: check for files in HEAD that are no longer in the index (i.e. deleted files).
    // for path := range treeFiles {
    //     if _, exists := index.GetEntry(path); !exists {
    //         stagedChanges = append(stagedChanges, fmt.Sprintf("deleted: %s", path))
    //     }
    // }

    return stagedChanges, nil
}
