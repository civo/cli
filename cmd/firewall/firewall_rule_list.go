package firewall

import (
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var firewallRuleListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Args:    cobra.MinimumNArgs(1),
	Example: "civo firewall rule ls FIREWALL_NAME",
	Short:   "List firewall rule",
	Long: `List all current firewall rules.
If you wish to use a custom format, the available fields are:

	* id
	* direction
	* action
	* protocol
	* start_port
	* end_port
	* cidr
	* label`,
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

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		firewallRules, err := client.ListFirewallRules(firewall.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, firewallRule := range firewallRules {
			ow.StartLine()

			if firewallRule.EndPort == "" && firewallRule.StartPort != "" {
				firewallRule.EndPort = firewallRule.StartPort
			}

			ow.AppendDataWithLabel("id", firewallRule.ID, "ID")
			// TODO: Check if is necessary this in the table, because you need pass like arg the name or the id of the firewall
			//ow.AppendDataWithLabel("Firewall", firewall.Name, "Firewall")
			ow.AppendDataWithLabel("direction", firewallRule.Direction, "Direction")
			ow.AppendDataWithLabel("protocol", firewallRule.Protocol, "Protocol")
			ow.AppendDataWithLabel("start_port", firewallRule.StartPort, "Start Port")
			ow.AppendDataWithLabel("end_port", firewallRule.EndPort, "End Port")
			ow.AppendDataWithLabel("action", firewallRule.Action, "Action")
			ow.AppendDataWithLabel("cidr", strings.Join(firewallRule.Cidr, ", "), "Cidr")
			ow.AppendDataWithLabel("label", firewallRule.Label, "Label")
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
