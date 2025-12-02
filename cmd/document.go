package cmd

import (
	"code-reviewer/internal/agent"
	"code-reviewer/internal/git"
	"code-reviewer/internal/ui"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var documentStaged bool
var documentUnstaged bool

var documentCmd = &cobra.Command{
	Use:   "document",
	Short: "Generate technical documentation for code changes",
	Run: func(cmd *cobra.Command, args []string) {
		if documentStaged && documentUnstaged {
			fmt.Fprintln(os.Stderr, "Error: --staged and --unstaged flags cannot be used together")
			os.Exit(1)
		}

		var mode git.DiffMode
		if documentStaged {
			mode = git.DiffModeStaged
		} else if documentUnstaged {
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
			fmt.Println("No changes to document.")
			return
		}

		a := agent.New()
		doc, err := ui.RunProgram(func() (string, error) {
			return a.Document(diff)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating documentation: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(doc)
	},
}

func init() {
	rootCmd.AddCommand(documentCmd)
	documentCmd.Flags().BoolVar(&documentStaged, "staged", false, "Document staged changes")
	documentCmd.Flags().BoolVar(&documentUnstaged, "unstaged", false, "Document unstaged changes")
}
