/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/app"
)

// addCmd represents the add command
var enableCmd = &cobra.Command{
	Use:     "enable",
	Aliases: []string{"add"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Enable or add one or more COPR repositories.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			repo, err := app.NewCoprRepo(arg)
			if err != nil {
				fmt.Println(err)
			}
			err = enableRepo(repo)
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

func verifyCoprRepo(r app.CoprRepo) error {
	_, err := http.Get(app.RepoFileUrl(r).String())
	if err != nil {
		return err
	}

	return nil
}

func enableRepo(r app.CoprRepo) error {
	if verifyCoprRepo(r) != nil {
		log.Fatal("Repository verification failed. Double-check repository name.")
	}
	if app.RepoExists(r) {
		err := app.ToggleRepo(r, app.Enabled)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := app.AddRepo(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(enableCmd)

}
