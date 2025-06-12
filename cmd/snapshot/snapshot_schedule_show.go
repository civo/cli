package snapshot

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var snapshotScheduleShowCmd = &cobra.Command{
	Use:     "show [ID/NAME]",
	Aliases: []string{"get", "inspect"},
	Short:   "Show snapshot schedule details",
	Long:    "Show detailed information about a specific snapshot schedule",
	Example: "civo snapshot schedule show my-schedule",
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

		data := map[string]string{
			"id":              schedule.ID,
			"name":            schedule.Name,
			"description":     schedule.Description,
			"cron_expression": schedule.CronExpression,
			"status":          schedule.Status.State,
			"paused":          strconv.FormatBool(schedule.Paused),
			"max_snapshots":   strconv.Itoa(schedule.Retention.MaxSnapshots),
			"created_at":      schedule.CreatedAt.Format(time.RFC822),
		}

		if schedule.Status.LastSnapshot.ID != "" {
			data["last_snapshot_id"] = schedule.Status.LastSnapshot.ID
			data["last_snapshot_name"] = schedule.Status.LastSnapshot.Name
			data["last_snapshot_state"] = schedule.Status.LastSnapshot.State
		}

		ow := utility.NewOutputWriterWithMap(data)

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow := utility.NewOutputWriter()
			ow.AppendDataWithLabel("ID", utility.Green(schedule.ID), "ID")
			ow.AppendDataWithLabel("Name", utility.Green(schedule.Name), "Name")
			ow.AppendDataWithLabel("Description", schedule.Description, "Description")
			ow.AppendDataWithLabel("Cron Expression", schedule.CronExpression, "Cron Expression")
			ow.AppendDataWithLabel("Status", utility.Green(schedule.Status.State), "Status")
			ow.AppendDataWithLabel("Paused", strconv.FormatBool(schedule.Paused), "Paused")
			ow.AppendDataWithLabel("Max Snapshots", strconv.Itoa(schedule.Retention.MaxSnapshots), "Max Snapshots")
			ow.AppendDataWithLabel("Created At", schedule.CreatedAt.Format(time.RFC822), "Created At")

			if schedule.Status.LastSnapshot.ID != "" {
				ow.AppendDataWithLabel("Last Snapshot ID", utility.Green(schedule.Status.LastSnapshot.ID), "Last Snapshot ID")
				ow.AppendDataWithLabel("Last Snapshot Name", utility.Green(schedule.Status.LastSnapshot.Name), "Last Snapshot Name")
				ow.AppendDataWithLabel("Last Snapshot State", utility.Green(schedule.Status.LastSnapshot.State), "Last Snapshot State")
			}

			ow.AppendDataWithLabel("Instances", "", "Instances")
			for i, instance := range schedule.Instances {
				k := fmt.Sprintf("Instance %d ID", i+1)
				ow.AppendDataWithLabel(k, utility.Green(instance.ID), k)
				if instance.Size != "" {
					k = fmt.Sprintf("Instance %d Size", i+1)
					ow.AppendDataWithLabel(k, instance.Size, k)
				}
				if len(instance.IncludedVolumes) > 0 {
					k = fmt.Sprintf("Instance %d Included Volumes", i+1)
					ow.AppendDataWithLabel(k, strings.Join(instance.IncludedVolumes, ", "), k)
				}
			}

			ow.WriteKeyValues()
		}
	},
}
