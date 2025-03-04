package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
)

// Computes and logs the difference between the working directory and the staging area
func Diff() error {
	index, err := LoadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %v", err)
	}
	wdMap, err := ScanWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to scan working directory: %v", err)
	}

	for _, entry := range *index {
		currentHash, exists := wdMap[entry.Path]
		if !exists {
			fmt.Printf("deleted: %s\n", entry.Path)
			continue
		}
		if currentHash != entry.BlobHash {
			currentContent, err := os.ReadFile(entry.Path)
			if err != nil {
				return fmt.Errorf("failed to read file: %v", err)
			}
			stagedContent, err := ReadBlobData(entry.BlobHash)
			if err != nil {
				return fmt.Errorf("failed to read blob data: %v", err)
			}
			diffs, err := ComputeDiff(string(currentContent), string(stagedContent), entry.Path)
			if err != nil {
				return fmt.Errorf("failed to compute diff: %v", err)
			}
			fmt.Print(diffs)
		}
	}
	return nil
}

func ComputeDiff(current, staged, path string) (string, error) {

	from := fmt.Sprintf("a/%s", path)
	to := fmt.Sprintf("b/%s", path)

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(staged),
		B:        difflib.SplitLines(current),
		FromFile: from, // Label for the staged version (older)
		ToFile:   to,   // Label for the working directory (newer)
		Context:  3,    // Number of context lines
	}

	diffText, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", err
	}
	coloredDiff := colorizeDiff(diffText)
	return coloredDiff, nil
}

func colorizeDiff(diffText string) string {
	var sb strings.Builder
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	lines := strings.Split(diffText, "\n")
	for _, line := range lines {
		// Do not colorize diff headers.
		if strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ") {
			sb.WriteString(line + "\n")
		} else if strings.HasPrefix(line, "@@"){
			sb.WriteString(blue(line) + "\n")
		} else if strings.HasPrefix(line, "-") {
			sb.WriteString(red(line) + "\n")
		} else if strings.HasPrefix(line, "+") {
			sb.WriteString(green(line) + "\n")
		} else {
			sb.WriteString(line + "\n")
		}
	}
	return sb.String()
}
