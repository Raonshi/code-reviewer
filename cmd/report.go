package cmd

import (
	"code-reviewer/internal/agent"
	"code-reviewer/internal/git"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a code review report",
	Run: func(cmd *cobra.Command, args []string) {
		// Try staged first
		diff, err := git.GetDiff(true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting staged diff: %v\n", err)
			os.Exit(1)
		}

		// If no staged changes, try unstaged
		if diff == "" {
			diff, err = git.GetDiff(false)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting unstaged diff: %v\n", err)
				os.Exit(1)
			}
		}

		if diff == "" {
			fmt.Println("No changes to review.")
			return
		}

		a := agent.New()
		report, err := a.Analyze(diff)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing code: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(report)
	},
}

func init() {
	reviewCmd.AddCommand(reportCmd)
}
