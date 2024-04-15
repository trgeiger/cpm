/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/app"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Remove one or more COPR repositories' configuration files.",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			repo, err := app.NewCoprRepo(arg)
			if err != nil {
				fmt.Println(err)
			}
			err = app.DeleteRepo(repo)
			if err != nil {
				if errors.Is(err, fs.ErrPermission) {
					fmt.Printf("This command must be run with superuser privileges.\nError: %s\n", err)
				} else {
					fmt.Println(err)
				}
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
