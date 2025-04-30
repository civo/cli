package snapshot

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var snapshotScheduleDeleteCmd = &cobra.Command{
	Use:     "delete [ID/NAME]",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete a snapshot schedule",
	Long:    "Delete a snapshot schedule by its ID or name",
	Example: "civo snapshot schedule delete my-schedule",
	Args:    cobra.ExactArgs(1),
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

		schedule, err := client.FindSnapshotSchedule(args[0])
		if err != nil {
			utility.Error("Finding snapshot schedule failed with %s", err)
			os.Exit(1)
		}

		_, err = client.DeleteSnapshotSchedule(schedule.ID)
		if err != nil {
			utility.Error("Deleting snapshot schedule failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{
			"result": "success",
			"id":     schedule.ID,
			"name":   schedule.Name,
		})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The snapshot schedule %s (%s) has been deleted\n", utility.Green(schedule.Name), utility.Green(schedule.ID))
		}
	},
}
