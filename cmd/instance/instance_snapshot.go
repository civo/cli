package instance

import (
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	instanceSnapshotName           string
	instanceSnapshotDescription    string
	instanceSnapshotIncludeVolumes bool
)

var instanceSnapshotCreateCmd = &cobra.Command{
	Use:     "snapshot create",
	Short:   "Create a new instance snapshot",
	Example: "civo instance snapshot create <instance-id> --name <snapshot-name>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		config := &civogo.CreateInstanceSnapshotConfig{
			Name:           instanceSnapshotName,
			Description:    instanceSnapshotDescription,
			IncludeVolumes: instanceSnapshotIncludeVolumes,
		}

		snapshot, err := client.CreateInstanceSnapshot(args[0], config)
		if err != nil {
			utility.Error("Creating instance snapshot failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendData("ID", snapshot.ID)
		ow.AppendData("Name", snapshot.Name)
		ow.AppendData("Description", snapshot.Description)
		ow.AppendData("Status", snapshot.Status.State)
		ow.AppendData("Created At", snapshot.CreatedAt.Format("2006-01-02 15:04:05"))
		ow.WriteTable()
	},
}
