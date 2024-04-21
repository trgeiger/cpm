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
	deleteAll bool
)

func NewRemoveCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"delete"},
		Args: func(cmd *cobra.Command, args []string) error {
			if deleteAll {
				if err := cobra.NoArgs(cmd, args); err != nil {
					return err
				}
			} else {
				if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
					return err
				}
			}
			return nil
		},
		Short: "Uninstall one or more Copr repositories.",
		Run: func(cmd *cobra.Command, args []string) {
			var repos []string
			if deleteAll {
				erepos, _ := app.GetReposList(fs, out, app.Enabled)
				drepos, _ := app.GetReposList(fs, out, app.Disabled)
				erepos = append(erepos, drepos...)
				for _, r := range erepos {
					repos = append(repos, r.Name())
				}
			} else {
				repos = args
			}
			for _, arg := range repos {
				repo, err := app.NewCoprRepo(arg)
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					err = app.DeleteRepo(repo, fs, out)
					if err != nil {
						app.SudoMessage(err, out)
					}
				}
			}
		},
	}
	cmd.Flags().BoolVarP(&deleteAll, "all", "A", false, "delete all installed Copr repositories")

	return cmd
}
