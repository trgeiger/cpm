/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/app"
)

func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed Copr repositories",
		Long: `A longer description that spans multiple lines and likely contains examples
				and usage of using your command. For example:
				
				Cobra is a CLI library for Go that empowers applications.
				This application is a tool to generate the needed files
				to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			repos, err := app.GetAllRepos()
			if err != nil {
				fmt.Printf("Error when retrieving locally installed repositories: %s", err)
			}
			showDupesMessage := false
			for _, r := range repos {
				fs := afero.NewOsFs()
				r.FindLocalFiles(fs)
				if len(r.LocalFiles) > 1 {
					showDupesMessage = true
				}
				fmt.Println(r.Name())
			}
			if showDupesMessage {
				fmt.Println("\nDuplicate entries found. Consider running the prune command.")
			}
		},
	}
}

func init() {
	rootCmd.AddCommand(NewListCmd())
}
