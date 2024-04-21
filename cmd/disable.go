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

func NewDisableCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "disable [repo(s)...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Disable one or more (space-separated) Copr repositories without uninstalling them.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				repo, err := app.NewCoprRepo(arg)
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					err = app.ToggleRepo(repo, fs, out, app.Disabled)
					if err != nil {
						app.SudoMessage(err, out)
					}
				}
			}
		},
	}
}
