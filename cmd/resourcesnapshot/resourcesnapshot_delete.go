package resourcesnapshot

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var resourceSnapshotDeleteCmd = &cobra.Command{
	Use:     "delete ID/NAME",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a resource snapshot",
	Long:    `Delete a resource snapshot by ID or name.`,
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

		if len(args) == 0 {
			utility.Error("No snapshot ID/NAME provided")
			os.Exit(1)
		}

		// Get the snapshot first to verify it exists and to show the name in output
		snapshot, err := client.GetResourceSnapshot(args[0])
		if err != nil {
			utility.Error("Error retrieving resource snapshot: %s", err)
			os.Exit(1)
		}

		// Delete the snapshot
		_, err = client.DeleteResourceSnapshot(snapshot.ID)
		if err != nil {
			utility.Error("Error deleting resource snapshot: %s", err)
			os.Exit(1)
		}

		if common.OutputFormat == "human" {
			fmt.Printf("The resource snapshot %s (%s) has been deleted\n", utility.Green(snapshot.Name), snapshot.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", snapshot.ID, "ID")
			ow.AppendDataWithLabel("name", snapshot.Name, "Name")
			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		}
	},
}
