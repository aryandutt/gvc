package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
	Parent  string
}

// CreateCommit creates a commit object from the current tree and updates HEAD.
func CreateCommit(message, author string) (string, error) {
	// 1. Create tree from index
	treeHash, err := CreateTreeFromIndex()
	if err != nil {
		return "", fmt.Errorf("create tree: %v", err)
	}

	// 2. Get parent commit (if exists)
	parentHash, err := getCurrentCommit() // Returns empty if no parent

	if err != nil {
		return "", fmt.Errorf("getting current commit: %v", err)
	}

	// 3. Build commit object
	commit := Commit{
		Hash:    "", // Will be set after creating the commit object
		Author:  author,
		Date:    time.Now(),
		Message: message,
		Parent:  parentHash,
	}

	// 4. Build commit content
	commitContent := fmt.Sprintf(
		"tree %s\n"+
			"parent %s\n"+
			"author %s %d +0000\n"+
			"committer %s %d +0000\n\n"+
			"%s\n",
		treeHash,
		commit.Parent,
		commit.Author,
		commit.Date.Unix(),
		commit.Author,
		commit.Date.Unix(),
		commit.Message,
	)

	// 5. Save commit object
	commitHash, err := CreateObject("commit", []byte(commitContent))
	if err != nil {
		return "", fmt.Errorf("create commit object: %v", err)
	}

	// 6. Update HEAD (current branch)
	if err := updateHead(commitHash); err != nil {
		return "", fmt.Errorf("update HEAD: %v", err)
	}

	commit.Hash = commitHash

	return commitHash, nil
}

func GetCommit(hash string) (*Commit, error) {
	// Read commit object from .gvc/objects
	objectPath := filepath.Join(".gvc", "objects", hash[:2], hash[2:])
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit %s: %v", hash, err)
	}

	// Split header and content
	parts := bytes.SplitN(data, []byte{0}, 2) // Git object format: "type size\0content"
	if len(parts) != 2 || !strings.HasPrefix(string(parts[0]), "commit") {
		return nil, fmt.Errorf("invalid commit object: %s", hash)
	}

	content := string(parts[1])
	lines := strings.Split(content, "\n")

	commit := &Commit{Hash: hash}
	for i, line := range lines {
		if strings.HasPrefix(line, "parent ") {
			commit.Parent = strings.TrimPrefix(line, "parent ")
		} else if strings.HasPrefix(line, "author ") {
			// Example: "author Alice <alice@example.com> 1700000000 +0000"
			fields := strings.Split(line, " ")
			if len(fields) < 4 {
				return nil, fmt.Errorf("invalid author details: %s", hash)
			}
			commit.Author = strings.Join(fields[1:len(fields)-2], " ") // Extract name/email

			// Parse timestamp
			timestamp, _ := strconv.ParseInt(fields[len(fields)-2], 10, 64)
			commit.Date = time.Unix(timestamp, 0)
		} else if line == "" {
			// Commit message starts after the first empty line
			commit.Message = strings.Join(lines[i+1:], "\n")
			break
		}
	}

	return commit, nil
}

// Helper: Get the current commit hash from HEAD
func getCurrentCommit() (string, error) {
	headRef, err := os.ReadFile(".gvc/HEAD")
	if err != nil {
		return "", err
	}

	// If HEAD points to a branch (e.g., "ref: refs/heads/main")
	refPath := string(headRef)[5:] // Remove "ref: "
	commitHash, err := os.ReadFile(filepath.Join(".gvc", refPath))
	if os.IsNotExist(err) {
		return "", nil // No parent (first commit)
	} else if err != nil {
		return "", err
	}

	return string(commitHash), nil
}

// Helper: Update HEAD (branch reference) to point to the new commit
func updateHead(commitHash string) error {
	headRef, err := os.ReadFile(".gvc/HEAD")
	if err != nil {
		return err
	}

	refPath := string(headRef)[5:] // Extract "refs/heads/main"
	return os.WriteFile(filepath.Join(".gvc", refPath), []byte(commitHash), 0644)
}
