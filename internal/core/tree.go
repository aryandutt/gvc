package core

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
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

	trees, err := createTreeRecursive("", index)
	if err != nil {
		return "", fmt.Errorf("failed to create tree: %v", err)
	}

	treeByte, err := ConvertTreeToByte(trees)
	if err != nil {
		return "", err
	}
	hash, err := CreateObject("tree", treeByte)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func createTreeRecursive(parent string, index *Index) ([]TreeEntry, error) {
	var tree []TreeEntry

	// Create a set to store children
	type void struct{}
	var member void
	childSet := make(map[string]void)

	// Iterate over index entries
	for _, entry := range *index {
		paths := strings.Split(entry.Path, "/")
		parents := strings.Split(parent, "/")

		// Remove common prefix
		for len(parents) > 0 && len(parent) > 0 {
			if parents[0] == paths[0] {
				parents = parents[1:]
				paths = paths[1:]
			} else {
				break;
			}
		}

		// Skip if the entry is not a child of the parent
		if len(paths) == 1 {
			if len(parents) != 0 {
				continue;
			}
			tree = append(tree, TreeEntry{
				Mode: entry.Type,
				Type: "blob",
				Hash: entry.BlobHash,
				Path: paths[0],
			})
			continue;
		}

		childSet[paths[0]] = member
	}
	
	// Recursively create child trees
	for child := range childSet {
		var childTrees []TreeEntry
		var err error
		if parent == "" {
			childTrees, err = createTreeRecursive(child, index) // [{0400 tree subdir_hash subdir} {100644 blob file_hash file}]
		} else {
			childTrees, err = createTreeRecursive(parent+"/"+child, index)
		}
		if err != nil {
			return nil, err
		}
		treeByte, err := ConvertTreeToByte(childTrees)
		if err != nil {
			return nil, err
		}
		hash, err := CreateObject("tree", treeByte)
		if err != nil {
			return nil, err
		}
		tree = append(tree, TreeEntry{
			Mode: "040000",
			Type: "tree",
			Hash: hash,
			Path: child,
		})
	}
	return tree, nil
}

// ConvertTreeToByte converts a list of tree entries to a byte slice.
func ConvertTreeToByte(tree []TreeEntry) ([]byte, error) {
	var treeContent bytes.Buffer
	for _, node := range tree {
		line := fmt.Sprintf("%s %s %s\t%s\n", node.Mode, node.Type, node.Hash, node.Path)
		treeContent.WriteString(line)
	}
	return treeContent.Bytes(), nil
}