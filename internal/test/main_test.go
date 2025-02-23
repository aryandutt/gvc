package test

import (
	"fmt"
	"os"
	"testing"
)

// write a basic test functio n
func TestFile(t *testing.T) {
	// Create a test directory

	info, _ := os.Lstat("../../dir/hello.txt")

	mode := info.Mode()

	fmt.Print((mode & 0111) != 0)
}
