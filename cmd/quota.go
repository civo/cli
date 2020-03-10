package cmd

import (
	"github.com/spf13/cobra"
)

var quotaCmd = &cobra.Command{
	Use:     "quota",
	Aliases: []string{"quotas"},
	Short:   "Your account's current quota settings and usage",
}

func init() {
	rootCmd.AddCommand(quotaCmd)
	quotaCmd.AddCommand(quotaShowCmd)
}
