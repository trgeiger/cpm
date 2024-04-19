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
	"github.com/trgeiger/copr-tool/internal/app"
)

func NewDisableCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Args:  cobra.MinimumNArgs(1),
		Short: "Disable one or more Copr repositories without uninstalling them.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				repo, err := app.NewCoprRepo(arg)
				if err != nil {
					fmt.Fprintln(out, err)
				}
				err = app.ToggleRepo(repo, fs, out, app.Disabled)
				if err != nil {
					app.SudoMessage(err, out)
					os.Exit(1)
				}
			}
		},
	}
}
