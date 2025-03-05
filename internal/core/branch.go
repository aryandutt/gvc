package core

import (
	"fmt"
	"os"
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
	newBranch := fmt.Sprintf(".gvc/refs/heads/%s", branch)
	_, err := os.Stat(newBranch)
	if os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist, use the -c flag to create a new branch", branch)
	} else if err != nil {
		return err
	}
	// Update HEAD to point to the new branch
	headRefPath := fmt.Sprintf("refs/heads/%s", branch)
	err = os.WriteFile(".gvc/HEAD", []byte(fmt.Sprintf("ref: %s", headRefPath)), 0644)
	return err
}
