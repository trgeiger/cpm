/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/trgeiger/cpm/internal/app"
)

func verifyCoprRepo(r *app.CoprRepo, fs afero.Fs) error {
	resp, err := http.Get(r.RepoUrl())
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return fmt.Errorf("repository does not exist, %s returned 404", r.RepoUrl())
	}
	resp, err = http.Get(r.RepoConfigUrl(fs))
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return fmt.Errorf("repository %s does not support Fedora release %s", r.Name(), app.FedoraReleaseVersion(fs))
	}
	return nil
}

func enableRepo(r *app.CoprRepo, fs afero.Fs, out io.Writer) error {
	if err := verifyCoprRepo(r, fs); err != nil {
		return err
	}
	err := r.FindLocalFiles(fs)
	if err != nil {
		return err
	}
	if r.LocalFileExists(fs) {
		err := app.ToggleRepo(r, fs, out, app.Enabled)
		if err != nil {
			app.SudoMessage(err, out)

			return err
		}
		return nil
	} else {
		err := app.AddRepo(r, fs, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewEnableCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:     "enable",
		Aliases: []string{"add"},
		Args:    cobra.MinimumNArgs(1),
		Short:   "Enable or add one or more Copr repositories.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				repo, err := app.NewCoprRepo(arg)
				if err != nil {
					fmt.Fprintln(out, err)
				} else {
					err = enableRepo(repo, fs, out)
					if err != nil {
						app.SudoMessage(err, out)
					}
				}
			}
		},
	}
}
