/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/app"
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
// 		fmt.Println("prune called")
// 	},
// }

func NewPruneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prune",
		Short: "Remove duplicate repository configurations.",
		Run: func(cmd *cobra.Command, args []string) {
			fs := afero.NewOsFs()
			repos, err := app.GetAllRepos()
			if err != nil {
				fmt.Printf("Error when retrieving locally installed repositories: %s", err)
				os.Exit(1)
			}
			pruneCount := 0
			for _, r := range repos {
				r.FindLocalFiles(fs)
				pruned, err := r.PruneDuplicates(fs)
				if pruned && err == nil {
					pruneCount++
				} else if pruned && err != nil {
					fmt.Printf("Pruning attempted on %s but encountered error: %s", r.Name(), err)
					os.Exit(1)
				} else if err != nil {
					fmt.Printf("Error encountered: %s", err)
					os.Exit(1)
				}
			}
			if pruneCount == 0 {
				fmt.Println("Nothing to prune.")
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewPruneCmd())
}
