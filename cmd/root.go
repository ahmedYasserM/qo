package cmd

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qo",
	Short: "A CLI tool to test students knowledge",
	Long:  "A CLI tool for managing Linux test challenges in a sandboxed environment.",
}

// Execute executes the root command.
func Execute() error {
	// Customize cobra output
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiYellow + cc.Bold,
		Commands:        cc.HiGreen + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.HiGreen + cc.Bold,
		Flags:           cc.HiBlue + cc.Bold,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	return rootCmd.Execute()
}
