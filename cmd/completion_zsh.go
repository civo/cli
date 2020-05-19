package cmd

import (
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// rootCmd.GenBashCompletion(os.Stdout)
		err := rootCmd.GenZshCompletion(os.Stdout)
		if err != nil {
			utility.Error("%s", err.Error())
			os.Exit(1)
		}
	},
}
