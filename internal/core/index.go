package core

import (
	"encoding/json"
	"os"
)

// IndexEntry represents a file in the staging area.
type IndexEntry struct {
	Path     string `json:"path"`     // File path relative to repo root
	BlobHash string `json:"blobHash"` // SHA-1 hash of the blob
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
