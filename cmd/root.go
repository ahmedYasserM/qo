package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eval",
	Short: "A CLI tool to test students knowledge",
	Long:  "A CLI tool for managing Linux test challenges in a sandboxed environment.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
