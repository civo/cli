package cmd

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var protocol, startPort, endPort, direction, label string
var cidr []string

var firewallRuleCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new firewall rule",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo firewall rule create FIREWALL_NAME or FIREWALL_ID [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			fmt.Printf("Unable to find the firewall for your search: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		newRuleConfig := &civogo.FirewallRuleConfig{
			FirewallID: firewall.ID,
			Protocol:   protocol,
			StartPort:  startPort,
			Cidr:       cidr,
			Direction:  direction,
			Label:      label,
		}

		if endPort == "" {
			newRuleConfig.EndPort = startPort
		} else {
			newRuleConfig.EndPort = endPort
		}

		rule, err := client.NewFirewallRule(newRuleConfig)
		if err != nil {
			fmt.Printf("Unable to create the new firewall rule: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": rule.ID, "Name": rule.Label})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("New rule %s created", aurora.Green(rule.Label))
		}
	},
}
