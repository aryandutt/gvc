package core

import (
	"fmt"

	"github.com/fatih/color"
)

func Status() error {
	index, err := LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %v", err)
	}

	wdMap, err := ScanWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to scan working directory: %v", err)
	}

	var modifiedFiles, deletedFiles, untrackedFiles []string

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
	stagedChanges, err := index.CompareToHead()
	if err != nil {
		return fmt.Errorf("failed to compare index to HEAD: %v", err)
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	if len(stagedChanges) > 0 {
		fmt.Println("\nChanges to be committed:")
		for _, change := range stagedChanges {
			fmt.Printf("\t%s\n", green(change))
		}
	}

	red := color.New(color.FgHiRed).SprintFunc()
	if len(modifiedFiles) > 0 || len(deletedFiles) > 0 {
		fmt.Println("\nChanges not staged for commit:")
		for _, path := range modifiedFiles {
			fmt.Printf("\t%s: %s\n", red("modified"), red(path))
		}
		for _, path := range deletedFiles {
			fmt.Printf("\t%s: %s\n", red("deleted"), red(path))
		}
	}

	if len(untrackedFiles) > 0 {
		fmt.Println("\nUntracked files:")
		for _, path := range untrackedFiles {
			fmt.Printf("\t%s: %s\n", red("untracked"), red(path))
		}
	}

	return nil
}
