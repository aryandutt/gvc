package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TreeEntry represents an entry in a tree object.
type TreeEntry struct {
	Mode string // e.g., "100644" for files
	Type string // "blob" or "tree"
	Hash string // SHA-1 hash of the blob/tree
	Path string // File or directory name
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

// GetHeadTree returns the hash of the tree object pointed by HEAD commit.
func GetHeadTree() (string, error) {
	commitHash, err := getCurrentCommit()
	if err != nil {
		return "", err
	}
	commit, err := GetCommit(commitHash)
	if err != nil {
		return "", err
	}
	return commit.Tree, nil
}

func GetTree(hash string) ([]TreeEntry, error) {
	objectPath := filepath.Join(".gvc", "objects", hash[:2], hash[2:])
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tree %s: %v", hash, err)
	}

	// Split header and content
	parts := bytes.SplitN(data, []byte{0}, 2) // Git object format: "type size\0content"
	if len(parts) != 2 || !strings.HasPrefix(string(parts[0]), "tree") {
		return nil, fmt.Errorf("invalid tree object: %s", hash)
	}

	content := string(parts[1])
	lines := strings.Split(content, "\n")

	tree := []TreeEntry{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		meta := strings.Split(parts[0], " ")
		tree = append(tree, TreeEntry{
			Mode: meta[0],
			Type: meta[1],
			Hash: meta[2],
			Path: parts[1],
		})
	}

	return tree, nil
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

		if parents[0] == "" {
			parents = parents[1:]
		}

		// Remove common prefix
		for len(parents) > 0 && len(parent) > 0 {
			if parents[0] == paths[0] {
				parents = parents[1:]
				paths = paths[1:]
			} else {
				break
			}
		}

		// Skip if the entry is not a child of the parent
		if len(paths) == 1 {
			if len(parents) != 0 {
				continue
			}
			tree = append(tree, TreeEntry{
				Mode: entry.Type,
				Type: "blob",
				Hash: entry.BlobHash,
				Path: paths[0],
			})
			continue
		}

		childSet[paths[0]] = member
	}

	// Recursively create child trees
	for child := range childSet {
		var childTrees []TreeEntry
		var err error
		if parent == "" {
			childTrees, err = createTreeRecursive(child, index)
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

// GetTree returns a map of file paths to their blob hash in a tree object.
func GetTreeFiles(hash string) (map[string]string, error) {
	tree, err := GetTree(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %v", err)
	}

	files := make(map[string]string)
	for _, entry := range tree {
		switch entry.Type {
		case "blob":
			files[entry.Path] = entry.Hash
		case "tree":
			// Recursively get files from the subtree.
			subFiles, err := GetTreeFiles(entry.Hash)
			if err != nil {
				return nil, fmt.Errorf("failed to get tree for %s: %v", entry.Path, err)
			}
			// Prepend the directory name to each file path from the sub-tree.
			for subPath, subHash := range subFiles {
				fullPath := filepath.Join(entry.Path, subPath)
				files[fullPath] = subHash
			}
		}
	}

	return files, nil
}
