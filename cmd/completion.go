package cmd

import (
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:     "completion [bash|zsh|powershell|fish]",
	Short:   "Generates bash completion scripts",
	Example: "civo completion [bash|zsh|powershell|fish]",
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	completionCmd.AddCommand(completionPowerShellCmd)
	completionCmd.AddCommand(completionFishCmd)
}
