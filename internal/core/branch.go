package core

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func SwitchBranch(branch string, create bool) error {
	// Read the HEAD file
	if create {
		currentCommit, err := getCurrentCommit()
		if err != nil {
			return err
		}
		newBranch := fmt.Sprintf(".gvc/refs/heads/%s", branch)
		err = os.WriteFile(newBranch, []byte(currentCommit), 0644)
		if err != nil {
			return err
		}
		// Update HEAD to point to the new branch
		headRefPath := fmt.Sprintf("refs/heads/%s", branch)
		err = os.WriteFile(".gvc/HEAD", []byte(fmt.Sprintf("ref: %s", headRefPath)), 0644)
		return err
	}
	isClean, err := IsWorkingDirClean()
	if err != nil {
		return err
	}
	if !isClean {
		return fmt.Errorf("working directory is not clean, commit or stash changes before switching branches")
	}
	newBranch := fmt.Sprintf(".gvc/refs/heads/%s", branch)
	_, err = os.Stat(newBranch)
	if os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist, use the -c flag to create a new branch", branch)
	} else if err != nil {
		return err
	}
	// Match the working dir with the new branch
	changedCommit, err := os.ReadFile(newBranch)
	if err != nil {
		return err
	}
	err = MatchDirectoryWithCommit(string(changedCommit))
	if err != nil {
		return err
	}
	// Update HEAD to point to the new branch
	headRefPath := fmt.Sprintf("refs/heads/%s", branch)
	err = os.WriteFile(".gvc/HEAD", []byte(fmt.Sprintf("ref: %s", headRefPath)), 0644)
	return err
}

func ListBranch() error {
	branches, err := os.ReadDir(".gvc/refs/heads")
	if err != nil {
		return err
	}
	// color the current branch green
	headRef, err := os.ReadFile(".gvc/HEAD")
	if err != nil {
		return err
	}
	green := color.New(color.FgGreen).SprintFunc()
	currentBranch := string(headRef)[16:]
	for _, branch := range branches {
		if branch.Name() == currentBranch {
			fmt.Printf("%s *\n", green(branch.Name()))
			continue
		}
		fmt.Println(branch.Name())
	}
	return nil
}
