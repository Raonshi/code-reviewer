package cmd

import (
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code changes",
	Long:  `Parent command for reviewing and fixing code changes.`,
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}
