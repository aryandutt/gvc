package core

import "os"

func InitRepo() error {
	// Create .gvc and subdirectories
	dirs := []string{".gvc", ".gvc/objects", ".gvc/refs", ".gvc/refs/heads"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	err := os.WriteFile(".gvc/index.json", []byte{}, 0644)
	if err != nil {
		return err
	}

	if err := os.WriteFile(".gvc/HEAD", []byte("ref: refs/heads/main"), 0644); err != nil {
		return err
	}

	return os.WriteFile(".gvc/refs/heads/main", []byte{}, 0644)
}
