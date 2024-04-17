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

func NewDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Args:  cobra.MinimumNArgs(1),
		Short: "Disable one or more Copr repositories without uninstalling them.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				repo, err := app.NewCoprRepo(arg)
				if err != nil {
					fmt.Println(err)
				}
				err = app.ToggleRepo(repo, app.Disabled)
				if err != nil {
					app.HandleError(err)
					os.Exit(1)
				}
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewDisableCmd())
}
