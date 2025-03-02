package cli

import (
  "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Use:   "gvc",
  Short: "A simple version control system written in Go",
  Long:  `GVC (Go Version Control) is a minimalist VCS for learning purposes.`,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    panic(err)
  }
}