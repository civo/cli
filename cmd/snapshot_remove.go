package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var snapshotRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Short:   "Remove a snapshot",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("snapshot") == nil {
			snapshot, err := client.FindSnapshot(args[0])
			if err != nil {
				utility.Error("Unable to find snapshot for your search %s", err)
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
				fmt.Printf("The snapshot called %s with ID %s was delete\n", utility.Green(snapshot.Name), utility.Green(snapshot.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
