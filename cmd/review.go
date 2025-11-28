package cmd

import (
	"fmt"
	"os"

	"code-reviewer/internal/agent"
	"code-reviewer/internal/git"

	"github.com/spf13/cobra"
)

var staged bool

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code changes",
	Long:  `Parent command for reviewing and fixing code changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsRepo() {
			fmt.Println("Error: Not a git repository.")
			os.Exit(1)
		}

		diff, err := git.GetDiff(staged)
		if err != nil {
			fmt.Printf("Error getting diff: %v\n", err)
			os.Exit(1)
		}

		if diff == "" {
			fmt.Println("No changes found.")
			return
		}

		fmt.Println("Analyzing changes...")
		ag := agent.New()
		report, err := ag.Analyze(diff)
		if err != nil {
			fmt.Printf("Error analyzing changes: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(report)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.Flags().BoolVar(&staged, "staged", false, "Review staged changes")
}
