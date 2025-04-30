package snapshot

import (
	"github.com/spf13/cobra"
)

// snapshotScheduleCmd represents the snapshot schedule command
var snapshotScheduleCmd = &cobra.Command{
	Use:     "schedule",
	Aliases: []string{"schedules"},
	Short:   "Manage snapshot schedules",
	Long:    `Create, list, update and delete snapshot schedules`,
}

func init() {
	snapshotScheduleCmd.AddCommand(snapshotScheduleCreateCmd)
	snapshotScheduleCmd.AddCommand(snapshotScheduleListCmd)
	snapshotScheduleCmd.AddCommand(snapshotScheduleShowCmd)
	snapshotScheduleCmd.AddCommand(snapshotScheduleUpdateCmd)
	snapshotScheduleCmd.AddCommand(snapshotScheduleDeleteCmd)
}
