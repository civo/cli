package resourcesnapshot

import (
	"fmt"
	"os"
	"time"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var resourceSnapshotShowCmd = &cobra.Command{
	Use:     "show ID/NAME",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"get", "inspect"},
	Short:   "Show details of a specific resource snapshot",
	Long: `Show details of a specific resource snapshot by ID or name.
If you wish to use a custom format, the available fields are:

* id
* name
* description
* resource_type
* created_at`,
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

		snapshot, err := client.GetResourceSnapshot(args[0])
		if err != nil {
			utility.Error("Error retrieving resource snapshot: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", snapshot.ID, "ID")
		ow.AppendDataWithLabel("name", snapshot.Name, "Name")
		ow.AppendDataWithLabel("description", snapshot.Description, "Description")
		ow.AppendDataWithLabel("resource_type", snapshot.ResourceType, "Resource Type")
		ow.AppendDataWithLabel("created_at", snapshot.CreatedAt.Format(time.RFC1123), "Created At")

		// Add instance-specific details if available
		if snapshot.Instance != nil {
			ow.AppendDataWithLabel("instance_id", snapshot.Instance.ID, "Instance ID")
			ow.AppendDataWithLabel("instance_name", snapshot.Instance.Name, "Instance Name")
			ow.AppendDataWithLabel("instance_status", snapshot.Instance.Status.State, "Instance Status")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if common.OutputFormat == "human" {
				fmt.Printf("Resource Snapshot: %s (%s)\n\n", utility.Green(snapshot.Name), snapshot.ID)
			}
			ow.WriteKeyValues()
		}
	},
}
