package cli

import (
	"fmt"

	"github.com/aryandutt/gvc/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
  }
  
  var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new gvc repository",
	Run: func(cmd *cobra.Command, args []string) {
	  if err := core.InitRepo(); err != nil {
		fmt.Println("Error:", err)
	  } else {
		fmt.Println("Initialized gvc repository")
	  }
	},
  }