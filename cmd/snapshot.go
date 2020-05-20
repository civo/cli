package cmd

import (
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Aliases: []string{"snapshots"},
	Short:   "Details of Civo snapshots",
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotRemoveCmd)

	snapshotCreateCmd.Flags().StringVarP(&cron, "cron", "c", "", "If a valid cron string is passed, the snapshot will be saved as an automated snapshot")
}
