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

var name, description string

var resourceSnapshotUpdateCmd = &cobra.Command{
	Use:     "update ID/NAME",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"modify", "edit"},
	Short:   "Update resource snapshot details",
	Long: `Update the name or description of a resource snapshot by ID or name.
If you wish to use a custom format, the available fields are:

* id
* name
* description`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if at least one flag is provided
		if name == "" && description == "" {
			return fmt.Errorf("you must provide at least one of --name or --description")
		}
		return nil
	},
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

		// Prepare update request
		req := &civogo.UpdateResourceSnapshotRequest{}
		if name != "" {
			req.Name = name
		}
		if description != "" {
			req.Description = description
		}

		// Update the snapshot
		updatedSnapshot, err := client.UpdateResourceSnapshot(snapshot.ID, req)
		if err != nil {
			utility.Error("Error updating resource snapshot: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", updatedSnapshot.ID, "ID")
		ow.AppendDataWithLabel("name", updatedSnapshot.Name, "Name")
		ow.AppendDataWithLabel("description", updatedSnapshot.Description, "Description")

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if common.OutputFormat == "human" {
				fmt.Printf("The resource snapshot %s (%s) has been updated\n", utility.Green(updatedSnapshot.Name), updatedSnapshot.ID)
			}
			ow.WriteKeyValues()
		}
	},
}

func init() {
	resourceSnapshotUpdateCmd.Flags().StringVarP(&name, "name", "n", "", "New name for the resource snapshot")
	resourceSnapshotUpdateCmd.Flags().StringVarP(&description, "description", "d", "", "New description for the resource snapshot")
}
