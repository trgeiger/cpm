/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/internal/app"
)

func NewListCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed Copr repositories",
		Long: `A longer description that spans multiple lines and likely contains examples
				and usage of using your command. For example:
				
				Cobra is a CLI library for Go that empowers applications.
				This application is a tool to generate the needed files
				to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			repos, err := app.GetAllRepos(fs, out)
			if err != nil {
				fmt.Fprintf(out, "Error when retrieving locally installed repositories: %s", err)
			}
			showDupesMessage := false
			for _, r := range repos {
				r.FindLocalFiles(fs)
				if len(r.LocalFiles) > 1 {
					showDupesMessage = true
				}
				fmt.Fprintln(out, r.Name())
			}
			if showDupesMessage {
				fmt.Fprintln(out, "\nDuplicate entries found. Consider running the prune command.")
			}
		},
	}
}
