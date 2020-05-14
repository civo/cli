package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// rootCmd.GenBashCompletion(os.Stdout)
		rootCmd.GenBashCompletion(os.Stdout)
	},
}
