package resourcesnapshot

import (
	"os"
	"time"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var resourceSnapshotListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all resource snapshots",
	Long: `List all resource snapshots currently available in your account.
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

		snapshots, err := client.ListResourceSnapshots()
		if err != nil {
			utility.Error("Error listing resource snapshots: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, snapshot := range snapshots {
			ow.StartLine()

			ow.AppendDataWithLabel("id", snapshot.ID, "ID")
			ow.AppendDataWithLabel("name", snapshot.Name, "Name")
			ow.AppendDataWithLabel("description", snapshot.Description, "Description")
			ow.AppendDataWithLabel("resource_type", snapshot.ResourceType, "Resource Type")
			ow.AppendDataWithLabel("created_at", snapshot.CreatedAt.Format(time.RFC1123), "Created At")
		}

		switch common.OutputFormat {
		case common.OutputFormatJSON:
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case common.OutputFormatCustom:
			ow.WriteCustomOutput(common.OutputFields)
		case common.OutputFormatHuman:
			ow.WriteTable()
		default:
			ow.WriteTable()
		}
	},
}
