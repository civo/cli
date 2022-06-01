package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo app show APP-NAME"`,
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show Application information",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		application, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("App %s", err)
			os.Exit(1)
		}

		networks, err := client.ListNetworks()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendDataWithLabel("id", application.ID, "ID")
		ow.AppendDataWithLabel("name", application.Name, "Name")

		for _, network := range networks {
			if application.NetworkID == network.ID {
				ow.AppendDataWithLabel("network_name", network.Name, "Network Name")
			}
		}

		ow.AppendDataWithLabel("region", client.Region, "Region")
		ow.AppendDataWithLabel("description", application.Description, "Description")
		ow.AppendDataWithLabel("image", application.Image, "Image")
		ow.AppendDataWithLabel("size", application.Size, "Size")
		ow.AppendDataWithLabel("domains", strings.Join(application.Domains, ", "), "Domains")
		ow.AppendDataWithLabel("status", application.Status, "Status")
		for _, process := range application.ProcessInfo {
			ow.AppendDataWithLabel("processType", process.ProcessType, "Process Type")
			ow.AppendDataWithLabel("processCount", strconv.Itoa(process.ProcessCount), "Process Count")
		}

		for _, config := range application.Config {
			//ow.AppendDataWithLabel("name", config.Name, "Config Name")
			ow.AppendDataWithLabel("value", config.Value, "Config Value")
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
