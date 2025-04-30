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

var scheduleDescription, scheduleCron string
var retentionPeriod string
var maxSnapshots int
var instanceIDs []string
var includeVolumes bool

var snapshotScheduleCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a new snapshot schedule",
	Long:    "Create a new snapshot schedule with specified cron expression and retention policy",
	Example: "civo snapshot schedule create --name my-schedule --cron '0 0 * * *' --instance-id instance-123 --retention-period 1w --max-snapshots 5",
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

		config := &civogo.CreateSnapshotScheduleRequest{
			Name:           scheduleName,
			Description:    scheduleDescription,
			CronExpression: scheduleCron,
			Retention: civogo.SnapshotRetention{
				Period:       retentionPeriod,
				MaxSnapshots: maxSnapshots,
			},
			Instances: make([]civogo.CreateSnapshotInstance, 0),
		}

		for _, instanceID := range instanceIDs {
			config.Instances = append(config.Instances, civogo.CreateSnapshotInstance{
				InstanceID:     instanceID,
				IncludeVolumes: includeVolumes,
			})
		}

		schedule, err := client.CreateSnapshotSchedule(config)
		if err != nil {
			utility.Error("Creating snapshot schedule failed with %s", err)
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
			fmt.Printf("Created snapshot schedule %s with ID %s\n", utility.Green(schedule.Name), utility.Green(schedule.ID))
		}
	},
}

func init() {
	snapshotScheduleCreateCmd.Flags().StringVarP(&scheduleName, "name", "n", "", "Name for the snapshot schedule")
	snapshotScheduleCreateCmd.Flags().StringVarP(&scheduleDescription, "description", "d", "", "Description for the snapshot schedule")
	snapshotScheduleCreateCmd.Flags().StringVarP(&scheduleCron, "cron", "c", "", "Cron expression for the schedule (e.g., '0 0 * * *' for daily at midnight)")
	snapshotScheduleCreateCmd.Flags().StringVarP(&retentionPeriod, "retention-period", "r", "", "Retention period for snapshots (e.g., '1w' for one week)")
	snapshotScheduleCreateCmd.Flags().IntVarP(&maxSnapshots, "max-snapshots", "m", 0, "Maximum number of snapshots to retain")
	snapshotScheduleCreateCmd.Flags().StringSliceVarP(&instanceIDs, "instance-id", "i", []string{}, "Instance IDs to snapshot (can be specified multiple times)")
	snapshotScheduleCreateCmd.Flags().BoolVarP(&includeVolumes, "include-volumes", "v", false, "Include attached volumes in snapshots")

	_ = snapshotScheduleCreateCmd.MarkFlagRequired("name")
	_ = snapshotScheduleCreateCmd.MarkFlagRequired("cron")
	_ = snapshotScheduleCreateCmd.MarkFlagRequired("instance-id")
}
