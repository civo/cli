package cmd

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var cron string

var snapshotCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new snapshot",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[1])
		if err != nil {
			utility.Error("Unable to find the instance %s", err)
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
			utility.Error("Unable to create the snapshot %s", err)
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
