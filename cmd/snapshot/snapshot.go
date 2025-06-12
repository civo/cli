package snapshot

import (
	"github.com/spf13/cobra"
)

// SnapshotCmd represents the snapshot command
var SnapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Aliases: []string{"snapshots"},
	Short:   "Manage snapshots and snapshot schedules",
	Long:    `Create, list, update and delete snapshots and snapshot schedules`,
}

func init() {
	SnapshotCmd.AddCommand(snapshotScheduleCmd)
}
