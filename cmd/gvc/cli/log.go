package cli

import (
	"fmt"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/fatih/color" // Optional: for colored output
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Display commit history",
	Run: func(cmd *cobra.Command, args []string) {
		commits, err := core.LogCommits()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for _, commit := range commits {
			yellow := color.New(color.FgYellow).SprintFunc()
			cyan := color.New(color.FgCyan).SprintFunc()
			fmt.Printf(
				"commit %s\nAuthor: %s\nDate:   %s\n\n    %s\n\n",
				yellow(commit.Hash[:7]),
				cyan(commit.Author),
				commit.Date.Format("Mon Jan 2 15:04:05 2006 -0700"),
				commit.Message,
			)
		}
	},
}
