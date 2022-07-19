package ip

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var ipListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo ip ls`,
	Short:   "List ips",
	Long:    `List all available reserved ips`,
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

		ips, err := client.ListIPs()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, ip := range ips.Items {
			ow.StartLine()
			ow.AppendDataWithLabel("id", ip.ID, "ID")
			ow.AppendDataWithLabel("name", ip.Name, "Name")
			ow.AppendDataWithLabel("address", ip.IP, "Address")
			if ip.AssignedTo.ID != "" {
				ow.AppendDataWithLabel("assigned_to", fmt.Sprintf("%s (%s)", ip.AssignedTo.Name, ip.AssignedTo.Type), "Assigned To(type)")
			} else {
				ow.AppendDataWithLabel("assigned_to", "No resource", "Assigned To(type)")
			}
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
