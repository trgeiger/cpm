/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/app"
)

func NewRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove",
		Aliases: []string{"delete"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Uninstall one or more Copr repositories.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				repo, err := app.NewCoprRepo(arg)
				if err != nil {
					fmt.Println(err)
				}
				err = app.DeleteRepo(repo)
				if err != nil {
					app.HandleError(err)
					os.Exit(1)
				}
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewRemoveCmd())
}
