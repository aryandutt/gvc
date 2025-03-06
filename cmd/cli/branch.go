package cli

import (
	"fmt"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)
}

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List all branches",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err:= core.ListBranch(); err != nil {
			fmt.Println("Error:", err)
		}
	},
}
