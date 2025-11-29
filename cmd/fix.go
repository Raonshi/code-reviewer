package cmd

import (
	"code-reviewer/internal/agent"
	"code-reviewer/internal/git"
	"code-reviewer/internal/ui"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Auto-fix code issues",
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
			fmt.Println("No changes to fix.")
			return
		}

		a := agent.New()
		fixedCode, err := ui.RunProgram(func() (string, error) {
			return a.Fix(diff)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating fix: %v\n", err)
			os.Exit(1)
		}

		// TODO: Apply fix to files
		// For now, we just print the fixed code
		fmt.Println("Proposed Fix:")
		fmt.Println(fixedCode)
		fmt.Println("Fix generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
