package snapshot

import (
	"os"
	"strconv"
	"time"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var snapshotScheduleListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all snapshot schedules",
	Long:    "List all snapshot schedules in your account",
	Example: "civo snapshot schedule list",
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

		schedules, err := client.ListSnapshotSchedules()
		if err != nil {
			utility.Error("Listing snapshot schedules failed with %s", err)
			os.Exit(1)
		}
		ow := utility.NewOutputWriter()

		switch common.OutputFormat {
		case "json":
			data := make([]map[string]string, 0)
			for _, schedule := range schedules {
				instanceCount := strconv.Itoa(len(schedule.Instances))
				lastSnapshot := "N/A"
				if schedule.Status.LastSnapshot.ID != "" {
					lastSnapshot = schedule.Status.LastSnapshot.State
				}

				data = append(data, map[string]string{
					"id":              schedule.ID,
					"name":            schedule.Name,
					"cron_expression": schedule.CronExpression,
					"status":          schedule.Status.State,
					"paused":          strconv.FormatBool(schedule.Paused),
					"instances":       instanceCount,
					"last_snapshot":   lastSnapshot,
					"created_at":      schedule.CreatedAt.Format(time.RFC822),
				})
			}
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if len(schedules) > 0 {
				for _, schedule := range schedules {
					ow.StartLine()
					ow.AppendDataWithLabel("id", utility.Green(schedule.ID), "ID")
					ow.AppendDataWithLabel("name", utility.Green(schedule.Name), "Name")
					ow.AppendDataWithLabel("cron_expression", schedule.CronExpression, "Cron Expression")
					ow.AppendDataWithLabel("status", utility.Green(schedule.Status.State), "Status")
					ow.AppendDataWithLabel("paused", strconv.FormatBool(schedule.Paused), "Paused")
					ow.AppendDataWithLabel("instances", strconv.Itoa(len(schedule.Instances)), "Instances")
					lastSnapshot := "N/A"
					if schedule.Status.LastSnapshot.ID != "" {
						lastSnapshot = schedule.Status.LastSnapshot.State
					}
					ow.AppendDataWithLabel("last_snapshot", lastSnapshot, "Last Snapshot")
					ow.AppendDataWithLabel("created_at", schedule.CreatedAt.Format(time.RFC822), "Created At")
				}
				ow.WriteTable()
			} else {
				utility.Info("No snapshot schedules found in %s.", client.Region)
			}
		}
	},
}
