package cmd

import (
	"code-reviewer/internal/git"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "code-reviewer",
	Short: "AI Code Review Agent CLI",
	Long:  `A CLI tool that uses AI to review code changes and suggest fixes.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !git.IsRepo() {
			fmt.Fprintln(os.Stderr, "Error: Current directory is not a git repository.")
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
