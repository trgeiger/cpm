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
	"github.com/spf13/viper"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fs := afero.NewOsFs()
	cmd, err := NewRootCmd(fs, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewRootCmd(fs afero.Fs, out io.Writer) (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:   "copr-tool",
		Short: "A command line tool for managing Copr repositories",
	}

	cmd.AddCommand(
		NewDisableCmd(fs, out),
		NewEnableCmd(fs, out),
		NewListCmd(fs, out),
		NewPruneCmd(fs, out),
		NewRemoveCmd(fs, out),
	)

	return cmd, nil
}

func init() {
	viper.SetConfigName("os-release")
	viper.SetConfigType("ini")
	viper.AddConfigPath("/etc/")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("could not fine /etc/os-release, copr-tool only functions on Fedora Linux systems: %w", err))

		} else {
			panic(fmt.Errorf("unknown fatal error: %w", err))
		}
	}
	if viper.Get("default.id") != "fedora" {
		fmt.Println("Non-Fedora distribution detected. Copr tool only functions on Fedora Linux.")
		os.Exit(1)
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.copr-tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
