/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/trgeiger/cpm/internal/app"
)

var (
	enabled  bool
	disabled bool
	showAll  bool
	verbose  bool
)

func printReposList(fs afero.Fs, out io.Writer, repoList []*app.CoprRepo) {
	if len(repoList) == 0 {
		return
	} else {
		showDupesMessage := false
		for _, r := range repoList {
			r.FindLocalFiles(fs)
			if len(r.LocalFiles) > 1 {
				showDupesMessage = true
			}
			fmt.Fprintln(out, r.Name())
		}
		if showDupesMessage {
			fmt.Fprintln(out, "\nDuplicate entries found. Consider running the prune command.")
		}
	}
}

func NewListCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed Copr repositories",
		Long: `List all installed and enabled Copr repositories by default.
For example:
	cpm list
	cpm list --all
	cpm list --disabled`,
		Run: func(cmd *cobra.Command, args []string) {
			var lists []app.RepoState
			if !disabled {
				lists = append(lists, app.Enabled)
			}
			if disabled || showAll {
				lists = append(lists, app.Disabled)
			}

			for i, list := range lists {
				var listState string
				if list == app.Disabled {
					listState = "disabled"
				} else {
					listState = "enabled"
				}
				printList, err := app.GetReposList(fs, out, list)
				if err != nil {
					fmt.Fprintf(out, "Error when retrieving local repositories: %s", err)
				}
				if len(printList) == 0 {
					fmt.Fprintf(out, "- No %s repositories\n", listState)
				} else {
					fmt.Fprintf(out, "- List of %s repositories:\n", listState)
					printReposList(fs, out, printList)
					if i == 0 && len(lists) > 1 {
						fmt.Fprintf(out, "\n")
					}
				}
			}

		},
	}

	cmd.Flags().BoolVarP(&enabled, "enabled", "e", true, "show enabled Copr repositories (default)")
	cmd.Flags().BoolVarP(&disabled, "disabled", "d", false, "show disabled Copr repositories")
	cmd.Flags().BoolVarP(&showAll, "all", "A", false, "show enabled and disabled Copr repositories")

	return cmd
}
