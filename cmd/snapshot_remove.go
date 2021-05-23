package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/civo/civogo"
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
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		snapshot, err := client.FindSnapshot(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s snapshot in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one snapshot with that value in your account")
				os.Exit(1)
			}
		}

		if utility.UserConfirmedDeletion("snapshot", defaultYes, snapshot.Name) {
			_, err = client.DeleteSnapshot(snapshot.Name)
			if err != nil {
				if errors.Is(err, civogo.DatabaseSnapshotCannotDeleteInUseError) {
					errMessage := fmt.Sprintf("sorry I couldn't delete this snapshot (%s) while it is in use\n", utility.Green(snapshot.Name))
					utility.Error(errMessage)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriterWithMap(map[string]string{"id": snapshot.ID, "name": snapshot.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON(prettySet)
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
