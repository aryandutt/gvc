package cli

import (
	"fmt"
	"log"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	switchCmd.Flags().BoolP("create", "c", false, "Create branch if it does not exist")
	rootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
	Use:   "switch [branch name]",
	Short: "Switch to a branch. Use -c flag to create the branch if it doesn't exist",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		createFlag, err := cmd.Flags().GetBool("create")
		if err != nil {
			log.Fatalf("failed to parse create flag: %v", err)
		}
		if err := core.SwitchBranch(branchName, createFlag); err != nil {
			log.Fatalf("failed to switch branch: %v", err)
		}
		fmt.Printf("Switched to branch %s\n", branchName)
	},
}
