package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:     "completion [bash|zsh|powershell|fish]",
	Short:   "Generates bash completion scripts",
	Example: "civo completion [bash|zsh|powershell|fish]",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	completionCmd.AddCommand(completionPowerShellCmd)
	completionCmd.AddCommand(completionFishCmd)
}
