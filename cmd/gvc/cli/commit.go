package cli

import (
	"fmt"
	"os/user"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/spf13/cobra"
)

var (
	message string
)

func init() {
	commitCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		if message == "" {
			fmt.Println("Error: Commit message required (-m)")
			return
		}
		
		user, err := user.Current()

		if err!=nil {
			fmt.Println("Error getting user info:", err)
			return
		}

		commitHash, err := core.CreateCommit(message, user.Username)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Committed: %s\n", commitHash[:7])
		}
	},
}
