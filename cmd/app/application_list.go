package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/common"
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

		if common.RegionSet != "" {
			client.Region = common.RegionSet
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

		ow := utility.NewOutputWriter()
		for _, app := range applications.Items {
			ow.StartLine()
			ow.AppendDataWithLabel("id", app.ID, "ID")
			ow.AppendDataWithLabel("name", app.Name, "Name")
			ow.AppendDataWithLabel("size", app.Size, "Size")
			ow.AppendDataWithLabel("network_id", app.NetworkID, "Network ID")
			ow.AppendDataWithLabel("domains", strings.Join(app.Domains, " "), "Domains")
			fmt.Println(app.ProcessInfo)
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
