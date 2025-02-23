package core

import (
	"bytes"
	"fmt"
	"sort"
)

// TreeEntry represents an entry in a tree object.
type TreeEntry struct {
	Mode string // e.g., "100644" for files
	Type string // "blob" or "tree"
	Hash string // SHA-1 hash of the blob/tree
	Path string // Filename or directory name
}

// CreateTreeFromIndex generates a tree object from the staging area (index.json).
func CreateTreeFromIndex() (string, error) {
	index, err := LoadIndex()
	if err != nil {
		return "", fmt.Errorf("failed to load index: %v", err)
	}

	// Sort entries for consistent hashing
	sort.Slice(*index, func(i, j int) bool {
		return (*index)[i].Path < (*index)[j].Path
	})

	// Build tree content
	var treeContent bytes.Buffer
	for _, entry := range *index {
		line := fmt.Sprintf("%s %s %s\t%s\n", "100644", "blob", entry.BlobHash, entry.Path)
		treeContent.WriteString(line)
	}

	// Create the tree object
	treeHash, err := CreateObject("tree", treeContent.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to create tree: %v", err)
	}

	return treeHash, nil
}
