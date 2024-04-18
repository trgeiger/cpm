/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/app"
)

func NewRemoveCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:     "remove",
		Aliases: []string{"delete"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Uninstall one or more Copr repositories.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				repo, err := app.NewCoprRepo(arg)
				if err != nil {
					fmt.Fprintln(out, err)
				}
				err = app.DeleteRepo(repo, fs, out)
				if err != nil {
					app.HandleError(err, out)
					os.Exit(1)
				}
			}
		},
	}
}
