package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/gorhill/cronexpr"
	"github.com/spf13/cobra"
)

var snapshotListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo snapshot ls`,
	Short:   "List snapshot",
	Long: `List all available snapshot.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* size_gigabytes
	* hostname
	* state
	* cron
	* schedule
	* schedule
	* requested_at
	* completed_at
	* instance_id
	* template
	* region
	* safe

Example: civo snapshot ls -o custom -f "id: name (hostname)"`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			return
		}

		snapshots, err := client.ListSnapshots()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		ow := utility.NewOutputWriter()

		for _, snapshot := range snapshots {
			ow.StartLine()
			ow.AppendDataWithLabel("id", snapshot.ID, "ID")
			ow.AppendDataWithLabel("name", snapshot.Name, "Name")
			ow.AppendDataWithLabel("size_gigabytes", fmt.Sprintf("%s GB", strconv.Itoa(snapshot.SizeGigabytes)), "Size")
			ow.AppendDataWithLabel("hostname", snapshot.Hostname, "Hostname")
			ow.AppendDataWithLabel("state", snapshot.State, "State")
			ow.AppendDataWithLabel("cron", snapshot.Cron, "Cron")
			if snapshot.Cron != "" {
				ow.AppendDataWithLabel("schedule", cronexpr.MustParse(snapshot.Cron).Next(time.Now()).Format(time.RFC1123), "Schedule")
			} else {
				ow.AppendDataWithLabel("schedule", "One-off", "Schedule")
			}
			ow.AppendDataWithLabel("requested_at", snapshot.RequestedAt.Format(time.RFC1123), "Requested At")
			ow.AppendDataWithLabel("completed_at", snapshot.CompletedAt.Format(time.RFC1123), "Completed At")

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendDataWithLabel("instance_id", snapshot.InstanceID, "InstanceID")
				ow.AppendDataWithLabel("template", snapshot.Template, "Template")
				ow.AppendDataWithLabel("region", snapshot.Region, "Region")
				ow.AppendDataWithLabel("safe", strconv.Itoa(snapshot.Safe), "Safe")
			}

		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
