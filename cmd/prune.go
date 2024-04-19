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

// pruneCmd represents the prune command
// var pruneCmd = &cobra.Command{
// 	Use:   "prune",
// 	Short: "A brief description of your command",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:

// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Fprintln(out, "prune called")
// 	},
// }

func NewPruneCmd(fs afero.Fs, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "prune",
		Short: "Remove duplicate repository configurations.",
		Run: func(cmd *cobra.Command, args []string) {
			repos, err := app.GetAllRepos(fs)
			if err != nil {
				fmt.Fprintf(out, "Error when retrieving locally installed repositories: %s", err)
				os.Exit(1)
			}
			pruneCount := 0
			for _, r := range repos {
				r.FindLocalFiles(fs)
				pruned, err := r.PruneDuplicates(fs, out)
				if pruned && err == nil {
					pruneCount++
				} else if pruned && err != nil {
					fmt.Fprintf(out, "Pruning attempted on %s but encountered error: %s", r.Name(), err)
					os.Exit(1)
				} else if err != nil {
					fmt.Fprintf(out, "Error encountered: %s", err)
					os.Exit(1)
				}
			}
			if pruneCount == 0 {
				fmt.Fprintln(out, "Nothing to prune.")
			}
		},
	}
}
