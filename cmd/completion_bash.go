package cmd

import (
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// rootCmd.GenBashCompletion(os.Stdout)
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			utility.Error("%s", err.Error())
			os.Exit(1)
		}
	},
}
