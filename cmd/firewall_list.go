package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var firewallListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List firewall",
	Long: `List all current firewalls.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* RulesCount
	* InstancesCount
	* Region

Example: civo firewall ls -o custom -f "ID: Name"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		firewalls, err := client.ListFirewalls()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, firewall := range firewalls {
			network, _ := client.FindNetwork(firewall.NetworkID)

			ow.StartLine()

			ow.AppendData("ID", firewall.ID)
			ow.AppendData("Name", firewall.Name)
			ow.AppendData("Network", network.Label)
			ow.AppendDataWithLabel("RulesCount", strconv.Itoa(firewall.RulesCount), "Total rules")
			ow.AppendDataWithLabel("InstancesCount", strconv.Itoa(firewall.InstancesCount), "Total Instances")
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
