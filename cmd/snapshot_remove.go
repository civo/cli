package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var snapshotRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo snapshot remove SNAPSHOT_NAME",
	Short:   "Remove/delete a snapshot",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion("snapshot", defaultYes) == true {
			snapshot, err := client.FindSnapshot(args[0])
			if err != nil {
				utility.Error("Finding snapshot for your search failed with %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteSnapshot(snapshot.Name)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": snapshot.ID, "Name": snapshot.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The snapshot called %s with ID %s was deleted\n", utility.Green(snapshot.Name), utility.Green(snapshot.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
