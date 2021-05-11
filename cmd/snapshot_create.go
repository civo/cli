package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var cron string

var snapshotCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo snapshot create SNAPSHOT_NAME INSTANCE_HOSTNAME",
	Short:   "Create a new snapshot",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[1])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		snapshotConfig := &civogo.SnapshotConfig{
			InstanceID: instance.ID,
			Safe:       false,
		}

		if cron != "" {
			snapshotConfig.Cron = cron
		}

		snapshot, err := client.CreateSnapshot(args[0], snapshotConfig)
		if err != nil {
			utility.Error("Creating the snapshot failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": snapshot.ID, "Name": snapshot.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a snapshot called %s with ID %s\n", utility.Green(snapshot.Name), utility.Green(snapshot.ID))
		}
	},
}
