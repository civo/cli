package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var firewallRuleListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "List firewall rule",
	Long: `List all current firewall rule.
If you wish to use a custom format, the available fields are:

	* ID
	* FirewallID
	* Direction
	* StartPort
	* EndPort
	* Label
	* Protocol
	* Cidr

Example: civo firewall rule ls -o custom -f "ID: Label"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			fmt.Printf("Unable to find firewall: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		firewallRules, err := client.ListFirewallRules(firewall.ID)
		if err != nil {
			fmt.Printf("Unable to list firewalls rules: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, firewallRule := range firewallRules {
			ow.StartLine()

			ow.AppendData("ID", firewallRule.ID)

			ow.AppendData("ID", firewallRule.ID)
			// TODO: Check if is necessary this in the table, because you need pass like arg the name or the id of the firewall
			//ow.AppendDataWithLabel("Firewall", firewall.Name, "Firewall")
			ow.AppendData("Direction", firewallRule.Direction)
			ow.AppendData("Protocol", firewallRule.Protocol)
			ow.AppendDataWithLabel("StartPort", firewallRule.StartPort, "Start Port")
			ow.AppendDataWithLabel("EndPort", firewallRule.EndPort, "End Port")
			ow.AppendData("Cidr", strings.Join(firewallRule.Cidr, ", "))
			ow.AppendData("Label", firewallRule.Label)
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
