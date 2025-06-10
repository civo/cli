package resourcesnapshot

import (
	"github.com/spf13/cobra"
)

// ResourceSnapshotCmd represents the resource-snapshot command
var ResourceSnapshotCmd = &cobra.Command{
	Use:     "resource-snapshot",
	Aliases: []string{"resourcesnapshot", "resource-snapshots", "resourcesnapshots"},
	Short:   "Manage resource snapshots",
	Long:    `List, update and delete resource snapshots`,
}

func init() {
	ResourceSnapshotCmd.AddCommand(resourceSnapshotListCmd)
	ResourceSnapshotCmd.AddCommand(resourceSnapshotShowCmd)
	ResourceSnapshotCmd.AddCommand(resourceSnapshotUpdateCmd)
	ResourceSnapshotCmd.AddCommand(resourceSnapshotDeleteCmd)
	ResourceSnapshotCmd.AddCommand(resourceSnapshotRestoreCmd)
}
