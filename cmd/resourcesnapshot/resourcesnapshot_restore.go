package resourcesnapshot

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	restoreDescription string
	restoreHostname    string
	restorePrivateIPv4 string
	overwriteExisting  bool
)

var resourceSnapshotRestoreCmd = &cobra.Command{
	Use:     "restore ID/NAME",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"recover"},
	Short:   "Restore from a resource snapshot",
	Long: `Restore a resource from a snapshot by ID or name.
If you wish to use a custom format, the available fields are:

* id
* name
* description
* resource_type`,
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

		// Get the snapshot first to verify it exists
		snapshot, err := client.GetResourceSnapshot(args[0])
		if err != nil {
			utility.Error("Error retrieving resource snapshot: %s", err)
			os.Exit(1)
		}

		// Prepare restore request
		req := &civogo.RestoreResourceSnapshotRequest{}

		// For now, we only support instance snapshots
		if snapshot.ResourceType == "instance" {
			instanceReq := &civogo.RestoreInstanceSnapshotRequest{
				Description:       restoreDescription,
				Hostname:          restoreHostname,
				PrivateIPv4:       restorePrivateIPv4,
				OverwriteExisting: overwriteExisting,
			}
			req.Instance = instanceReq
		} else {
			utility.Error("Unsupported snapshot type: %s. Currently only instance snapshots are supported", snapshot.ResourceType)
			os.Exit(1)
		}

		// Restore from the snapshot
		restoredSnapshot, err := client.RestoreResourceSnapshot(snapshot.ID, req)
		if err != nil {
			utility.Error("Error restoring from resource snapshot: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("resource_type", restoredSnapshot.ResourceType, "Resource Type")

		if restoredSnapshot.ResourceType == "instance" && restoredSnapshot.Instance != nil {
			ow.AppendDataWithLabel("id", restoredSnapshot.Instance.ID, "ID")
			ow.AppendDataWithLabel("name", restoredSnapshot.Instance.Name, "Name")
			ow.AppendDataWithLabel("description", restoredSnapshot.Instance.Description, "Description")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if common.OutputFormat == common.OutputFormatHuman {
				fmt.Printf("The resource of type %s has been restored from snapshot %s (%s).\n",
					utility.Green(restoredSnapshot.ResourceType),
					utility.Green(snapshot.Name),
					snapshot.ID)
			}
			ow.WriteKeyValues()
		}
	},
}

func init() {
	resourceSnapshotRestoreCmd.Flags().StringVarP(&restoreDescription, "description", "d", "", "Description for the restored resource")
	resourceSnapshotRestoreCmd.Flags().StringVarP(&restoreHostname, "hostname", "n", "", "Hostname for the restored instance (instance snapshots only)")
	resourceSnapshotRestoreCmd.Flags().StringVar(&restorePrivateIPv4, "private-ip", "", "Private IP for the restored instance (instance snapshots only)")
	resourceSnapshotRestoreCmd.Flags().BoolVar(&overwriteExisting, "overwrite", false, "Overwrite existing instance if it exists (instance snapshots only)")
}
