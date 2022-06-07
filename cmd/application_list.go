package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List all appplications",
	Example: "civo app ls",
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

		applications, err := client.ListApplications()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		networks, err := client.ListNetworks()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		for _, app := range applications.Items {
			ow.AppendDataWithLabel("id", app.ID, "ID")
			ow.AppendDataWithLabel("name", app.Name, "Name")
			ow.AppendDataWithLabel("size", app.Size, "Size")
			for _, network := range networks {
				if app.NetworkID == network.ID {
					ow.AppendDataWithLabel("network_name", network.Label, "Network Name")
				}
			}
			ow.AppendDataWithLabel("status", app.Status, "Status")
			if len(app.Domains) > 1 {
				ow.AppendDataWithLabel("domains", fmt.Sprintf(app.Domains[0]+" ..."), "Domains")
			} else {
				ow.AppendDataWithLabel("domains", app.Domains[0], "Domains")
			}

			for _, process := range app.ProcessInfo {
				ow.AppendDataWithLabel("process_type", process.ProcessType, "Process Type")
				ow.AppendDataWithLabel("process_count", strconv.Itoa(process.ProcessCount), "Process Count")
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
