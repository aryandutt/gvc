package cli

import (
	"fmt"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(diffCmd)
  }
  
  var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Shows the difference between the working directory and the staging area",
	Run: func(cmd *cobra.Command, args []string) {
		if err := core.Diff(); err != nil {
			fmt.Println("Error:", err)
		}
	},
  }