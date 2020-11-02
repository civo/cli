package cmd

import (
	"os"

	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var completionFishCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generates fish completion scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// rootCmd.GenBashCompletion(os.Stdout)
		err := rootCmd.GenFishCompletion(os.Stdout, true)
		if err != nil {
			utility.Error("%s", err.Error())
			os.Exit(1)
		}
	},
}
