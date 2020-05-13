package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var firewallListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List firewall",
	Long: `List all current firewall.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* RulesCount
	* InstancesCount
	* Region

Example: civo firewall ls -o custom -f "ID: Name"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		firewalls, err := client.ListFirewalls()
		if err != nil {
			utility.Error("Unable to list firewalls %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, firewall := range firewalls {
			ow.StartLine()

			ow.AppendData("ID", firewall.ID)
			ow.AppendData("Name", firewall.Name)
			ow.AppendDataWithLabel("RulesCount", strconv.Itoa(firewall.RulesCount), "Total rules")
			ow.AppendDataWithLabel("InstancesCount", strconv.Itoa(firewall.InstancesCount), "Total Intances")
			ow.AppendData("Region", firewall.Region)
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
