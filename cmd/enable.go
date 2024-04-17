/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trgeiger/copr-tool/app"
)

func verifyCoprRepo(r *app.CoprRepo) error {
	resp, err := http.Get(r.RepoUrl())
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return fmt.Errorf("repository does not exist, %s returned 404", r.RepoUrl())
	}
	resp, err = http.Get(r.RepoConfigUrl())
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return fmt.Errorf("repository %s does not support Fedora release %s", r.Name(), app.FedoraReleaseVersion())
	}
	return nil
}

func enableRepo(r *app.CoprRepo) error {
	if err := verifyCoprRepo(r); err != nil {
		return err
	}
	fs := afero.NewOsFs()
	err := r.FindLocalFiles(fs)
	if err != nil {
		return err
	}
	// if len(r.LocalFiles) > 1 {
	// 	err := r.PruneDuplicates(fs)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	if r.LocalFileExists() {
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

func NewEnableCmd(config *viper.Viper) *cobra.Command {
	return &cobra.Command{
		Use:     "enable",
		Aliases: []string{"add"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Enable or add one or more Copr repositories.",
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
					app.HandleError(err)
					os.Exit(1)
				}
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewEnableCmd(viper.GetViper()))

}
