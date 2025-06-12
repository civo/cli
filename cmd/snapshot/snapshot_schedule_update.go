package snapshot

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var scheduleName string
var schedulePaused string

var snapshotScheduleUpdateCmd = &cobra.Command{
	Use:     "update [ID/NAME]",
	Short:   "Update a snapshot schedule",
	Long:    "Update the properties of an existing snapshot schedule",
	Example: "civo snapshot schedule update my-schedule --name new-name --description 'New description' --pause",
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

		var pauseSchedule *bool
		switch schedulePaused {
		case "true", "yes":
			pauseSchedule = utility.Ptr(true)
		case "false", "no":
			pauseSchedule = utility.Ptr(false)
		case "":
			pauseSchedule = nil
		default:
			utility.Error("Invalid value for -paused, please use true/false values.")
			os.Exit(1)
		}

		config := &civogo.UpdateSnapshotScheduleRequest{
			Name:        scheduleName,
			Description: scheduleDescription,
			Paused:      pauseSchedule,
		}

		schedule, err := client.UpdateSnapshotSchedule(args[0], config)
		if err != nil {
			utility.Error("Updating snapshot schedule failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{
			"id":   schedule.ID,
			"name": schedule.Name,
		})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Updated snapshot schedule %s\n", utility.Green(schedule.Name))
		}
	},
}

func init() {
	snapshotScheduleUpdateCmd.Flags().StringVarP(&scheduleName, "name", "n", "", "New name for the snapshot schedule")
	snapshotScheduleUpdateCmd.Flags().StringVarP(&scheduleDescription, "description", "d", "", "New description for the snapshot schedule")
	snapshotScheduleUpdateCmd.Flags().StringVarP(&schedulePaused, "paused", "p", "", "Whether to pause the snapshot schedule (use 'true' or 'false')")
}
