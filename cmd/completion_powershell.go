package cmd

import (
	"os"

	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var completionPowerShellCmd = &cobra.Command{
	Use:   "powershell",
	Short: "Generates powershell completion scripts",
	Run: func(cmd *cobra.Command, args []string) {
		// rootCmd.GenBashCompletion(os.Stdout)
		err := rootCmd.GenPowerShellCompletion(os.Stdout)
		if err != nil {
			utility.Error("%s", err.Error())
			os.Exit(1)
		}
	},
}
