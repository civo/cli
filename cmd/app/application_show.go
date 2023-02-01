package app

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "info"},
	Example: `civo app show APP_NAME`,
	Short:   "Prints information about an App",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		app, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		ow.StartLine()
		fmt.Println()
		ow.AppendDataWithLabel("id", app.ID, "ID")
		ow.AppendDataWithLabel("name", app.Name, "Name")
		ow.AppendDataWithLabel("size", app.Size, "Size")
		ow.AppendDataWithLabel("network_id", app.NetworkID, "Network ID")
		ow.AppendDataWithLabel("firewall_id", app.FirewallID, "Firewall ID")
		ow.AppendDataWithLabel("app_ip", app.AppIP, "App IP")
		ow.AppendDataWithLabel("status", app.Status, "Status")
		if app.GitInfo != nil {
			ow.AppendDataWithLabel("git_url", app.GitInfo.GitURL, "Git URL")
		}

		// TODO: Separate this into a separate command
		// if app.ProcessInfo != nil {
		// 	for _, process := range app.ProcessInfo {
		// 		ow.AppendDataWithLabel("process_type", process.ProcessType, "Process Type")
		// 		ow.AppendDataWithLabel("process_count", strconv.Itoa(process.ProcessCount), "Process Count")
		// 	}
		// }

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
