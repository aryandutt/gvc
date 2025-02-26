package test

import (
	"fmt"
	"testing"
)

// write a basic test functio n
func TestFile(t *testing.T) {
	// Create a test directory

	arr := []string{"file1.txt", "file2.txt"}
	arr = arr[1:]

	fmt.Print(arr)
}
