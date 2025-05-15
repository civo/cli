package instance

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:     "snapshot",
	Aliases: []string{"snapshots"},
	Short:   "Manage instance snapshots",
}

var snapshotCreateCmd = &cobra.Command{
	Use:     "create [INSTANCE_NAME/ID]",
	Short:   "Create a snapshot of an instance",
	Example: "civo instance snapshot create INSTANCE_NAME/ID --name snapshot-name",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		config := &civogo.CreateInstanceSnapshotParams{
			Name:        name,
			Description: description,
		}

		snapshot, err := client.CreateInstanceSnapshot(args[0], config)
		if err != nil {
			utility.Error("Creating snapshot failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", snapshot.ID, "ID")
		ow.AppendDataWithLabel("name", snapshot.Name, "Name")
		ow.AppendDataWithLabel("description", snapshot.Description, "Description")
		ow.AppendDataWithLabel("state", snapshot.Status.State, "State")
		ow.AppendDataWithLabel("created_at", snapshot.CreatedAt.Format("2006-01-02 15:04:05"), "Created At")

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}

var snapshotListCmd = &cobra.Command{
	Use:     "list [INSTANCE_NAME/ID]",
	Short:   "List all snapshots of an instance",
	Example: "civo instance snapshot list INSTANCE_NAME/ID",
	Aliases: []string{"ls"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		snapshots, err := client.ListInstanceSnapshots(args[0])
		if err != nil {
			utility.Error("Listing snapshots failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, snapshot := range snapshots {
			ow.StartLine()
			ow.AppendDataWithLabel("id", snapshot.ID, "ID")
			ow.AppendDataWithLabel("name", snapshot.Name, "Name")
			ow.AppendDataWithLabel("description", snapshot.Description, "Description")
			ow.AppendDataWithLabel("state", snapshot.Status.State, "State")
			ow.AppendDataWithLabel("created_at", snapshot.CreatedAt.Format("2006-01-02 15:04:05"), "Created At")
		}

		// Use FinishAndPrintOutput to handle empty results and different output formats
		ow.FinishAndPrintOutput()
	},
}

var snapshotShowCmd = &cobra.Command{
	Use:     "show [INSTANCE_NAME/ID] [SNAPSHOT_NAME/ID]",
	Short:   "Show details of an instance snapshot",
	Example: "civo instance snapshot show INSTANCE_NAME/ID SNAPSHOT_NAME/ID",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		snapshot, err := client.GetInstanceSnapshot(args[0], args[1])
		if err != nil {
			utility.Error("Getting snapshot failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", snapshot.ID, "ID")
		ow.AppendDataWithLabel("name", snapshot.Name, "Name")
		ow.AppendDataWithLabel("description", snapshot.Description, "Description")
		ow.AppendDataWithLabel("state", snapshot.Status.State, "State")
		ow.AppendDataWithLabel("created_at", snapshot.CreatedAt.Format("2006-01-02 15:04:05"), "Created At")

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}

var snapshotUpdateCmd = &cobra.Command{
	Use:     "update [INSTANCE_NAME/ID] [SNAPSHOT_NAME/ID]",
	Short:   "Update an instance snapshot",
	Example: "civo instance snapshot update INSTANCE_NAME/ID SNAPSHOT_NAME/ID --name new-name --description new-description",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		params := &civogo.UpdateInstanceSnapshotParams{
			Name:        name,
			Description: description,
		}

		snapshot, err := client.UpdateInstanceSnapshot(args[0], args[1], params)
		if err != nil {
			utility.Error("Updating snapshot failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", snapshot.ID, "ID")
		ow.AppendDataWithLabel("name", snapshot.Name, "Name")
		ow.AppendDataWithLabel("description", snapshot.Description, "Description")
		ow.AppendDataWithLabel("state", snapshot.Status.State, "State")
		ow.AppendDataWithLabel("created_at", snapshot.CreatedAt.Format("2006-01-02 15:04:05"), "Created At")

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:     "delete [INSTANCE_NAME/ID] [SNAPSHOT_NAME/ID]",
	Short:   "Delete an instance snapshot",
	Example: "civo instance snapshot delete INSTANCE_NAME/ID SNAPSHOT_NAME/ID",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		err = client.DeleteInstanceSnapshot(args[0], args[1])
		if err != nil {
			utility.Error("Deleting snapshot failed with %s", err)
			os.Exit(1)
		}

		fmt.Println(fmt.Sprintf("The snapshot %s has been deleted", utility.Green(args[1])))
	},
}

// Added: Flag variables for snapshot restore
var (
	restoreDescription       string
	restoreHostname          string
	restorePrivateIPv4       string
	restoreOverwriteExisting bool
)

var snapshotRestoreCmd = &cobra.Command{
	Use:     "restore [INSTANCE_NAME/ID] [SNAPSHOT_NAME/ID]",
	Short:   "Restore an instance from a snapshot",
	Example: "civo instance snapshot restore INSTANCE_NAME/ID SNAPSHOT_NAME/ID",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		// Construct params from flags
		params := civogo.RestoreInstanceSnapshotParams{
			Description:       restoreDescription,
			Hostname:          restoreHostname,
			PrivateIPv4:       restorePrivateIPv4,
			OverwriteExisting: restoreOverwriteExisting,
		}

		// Updated SDK call to RestoreInstanceSnapshot
		instanceRestoreInfo, err := client.RestoreInstanceSnapshot(args[0], args[1], &params)
		if err != nil {
			utility.Error("Restoring snapshot failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", instanceRestoreInfo.ID, "ID")
		ow.AppendDataWithLabel("name", instanceRestoreInfo.Name, "Name")
		ow.AppendDataWithLabel("hostname", instanceRestoreInfo.Hostname, "Hostname")
		ow.AppendDataWithLabel("description", instanceRestoreInfo.Description, "Description")
		ow.AppendDataWithLabel("from_snapshot", instanceRestoreInfo.FromSnapshot, "From Snapshot")
		ow.AppendDataWithLabel("private_ipv4", instanceRestoreInfo.PrivateIPv4, "Private IPv4")
		ow.AppendDataWithLabel("overwrite_existing", fmt.Sprintf("%t", instanceRestoreInfo.OverwriteExisting), "Overwrite Existing")
		ow.AppendDataWithLabel("state", instanceRestoreInfo.Status.State, "State")
		ow.AppendDataWithLabel("created_at", instanceRestoreInfo.CreatedAt.Format("2006-01-02 15:04:05"), "Created At")
		if instanceRestoreInfo.CompletedAt != nil {
			ow.AppendDataWithLabel("completed_at", instanceRestoreInfo.CompletedAt.Format("2006-01-02 15:04:05"), "Completed At")
		} else {
			ow.AppendDataWithLabel("completed_at", "", "Completed At")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			// Print a success message before the key-value output for default format
			fmt.Printf("The instance %s has been restored from snapshot %s. Details:\n", utility.Green(args[0]), utility.Green(args[1]))
			ow.WriteKeyValues()
		}
	},
}

func init() {
	snapshotCreateCmd.Flags().StringP("name", "n", "", "Name for the snapshot")
	snapshotCreateCmd.Flags().StringP("description", "d", "", "Description for the snapshot")
	snapshotCreateCmd.MarkFlagRequired("name")

	snapshotUpdateCmd.Flags().StringP("name", "n", "", "New name for the snapshot")
	snapshotUpdateCmd.Flags().StringP("description", "d", "", "New description for the snapshot")

	snapshotRestoreCmd.Flags().StringVarP(&restoreDescription, "description", "d", "", "New description for the restore operation")
	snapshotRestoreCmd.Flags().StringVar(&restoreHostname, "hostname", "", "New hostname for the restored instance (optional)")
	snapshotRestoreCmd.Flags().StringVar(&restorePrivateIPv4, "private-ipv4", "", "New private IPv4 address for the restored instance (optional)")
	snapshotRestoreCmd.Flags().BoolVar(&restoreOverwriteExisting, "overwrite-existing", false, "Overwrite the existing instance")

	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotShowCmd)
	snapshotCmd.AddCommand(snapshotUpdateCmd)
	snapshotCmd.AddCommand(snapshotDeleteCmd)
	snapshotCmd.AddCommand(snapshotRestoreCmd)
}
