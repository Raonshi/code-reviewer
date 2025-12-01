package cmd

import (
	"code-reviewer/internal/agent"
	"code-reviewer/internal/git"
	"code-reviewer/internal/ui"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var reportStaged bool
var reportUnstaged bool

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a code review report",
	Run: func(cmd *cobra.Command, args []string) {
		if reportStaged && reportUnstaged {
			fmt.Fprintln(os.Stderr, "Error: --staged and --unstaged flags cannot be used together")
			os.Exit(1)
		}

		var mode git.DiffMode
		if reportStaged {
			mode = git.DiffModeStaged
		} else if reportUnstaged {
			mode = git.DiffModeUnstaged
		} else {
			mode = git.DiffModeAll
		}

		diff, err := git.GetDiff(mode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting diff: %v\n", err)
			os.Exit(1)
		}

		if diff == "" {
			fmt.Println("No changes to review.")
			return
		}

		a := agent.New()
		report, err := ui.RunProgram(func() (string, error) {
			return a.Analyze(diff)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing code: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(report)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().BoolVar(&reportStaged, "staged", false, "Review staged changes")
	reportCmd.Flags().BoolVar(&reportUnstaged, "unstaged", false, "Review unstaged changes")
}
