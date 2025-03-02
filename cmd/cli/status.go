package cli

import (
	"fmt"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows status of the repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := core.Status(); err != nil {
			fmt.Println("Error:", err)
		}
	},
}
