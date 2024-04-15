/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trgeiger/copr-tool/app"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		listRepos()
	},
}

func listRepos() error {
	files, err := os.ReadDir(app.ReposDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			ioFile, err := os.Open(app.ReposDir + file.Name())

			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(ioFile)
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "[copr:copr") {
					t := strings.Split(strings.Trim(scanner.Text(), "[]"), ":")
					r, _ := app.NewCoprRepo(t[len(t)-2] + "/" + t[len(t)-1])
					properFileName := app.RepoFileName(r)
					if file.Name() != properFileName {
						fmt.Printf("Repository %s detected with non-standard repository file name.", r.User+"/"+r.Project)
					}
					fmt.Println(strings.Join(t[len(t)-2:], "/"))
					break
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Issue reading repo files: ", err)
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
}
