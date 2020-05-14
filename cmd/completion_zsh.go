package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// rootCmd.GenBashCompletion(os.Stdout)
		rootCmd.GenZshCompletion(os.Stdout)
	},
}
