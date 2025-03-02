package cli

import (
	"fmt"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [file]",
	Short: "Add files to the staging area",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			if err := core.AddToStage(file); err != nil {
				fmt.Printf("Error adding %s: %v\n", file, err)
			} else {
				fmt.Printf("Added %s\n", file)
			}
		}
	},
}
