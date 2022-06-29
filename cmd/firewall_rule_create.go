package cmd

import (
	"fmt"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var protocol, startPort, endPort, direction, label, cidr, action, directionValue string

var firewallRuleCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new firewall rule",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo firewall rule create FIREWALL_NAME/FIREWALL_ID [flags]",
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

		newRuleConfig := &civogo.FirewallRuleConfig{
			FirewallID: firewall.ID,
			Protocol:   protocol,
			StartPort:  startPort,
			Cidr:       strings.Split(cidr, ","),
			Label:      label,
			Action:     action,
		}

		// Check the rule address, if the input is different
		// from (ingress or egress) then we will generate an error
		if direction == "ingress" {
			newRuleConfig.Direction = direction
			directionValue = "from"
		} else if direction == "egress" {
			newRuleConfig.Direction = direction
			directionValue = "to"
		} else {
			utility.Error("'--direction' flag can't be empty")
			os.Exit(1)
		}

		if endPort == "" {
			newRuleConfig.EndPort = startPort
		} else {
			newRuleConfig.EndPort = endPort
		}

		rule, err := client.NewFirewallRule(newRuleConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": rule.ID, "name": rule.Label})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if rule.Label == "" {
				if newRuleConfig.EndPort == newRuleConfig.StartPort {
					fmt.Printf("Created a firewall rule to %s, %s access to port %s %s %s with ID %s\n", utility.Green(newRuleConfig.Action), utility.Green(newRuleConfig.Direction), utility.Green(newRuleConfig.StartPort), directionValue, utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				} else {
					fmt.Printf("Created a firewall rule to %s, %s access to ports %s-%s %s %s with ID %s\n", utility.Green(newRuleConfig.Action), utility.Green(newRuleConfig.Direction), utility.Green(newRuleConfig.StartPort), utility.Green(newRuleConfig.EndPort), directionValue, utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				}
			} else {
				if newRuleConfig.EndPort == newRuleConfig.StartPort {
					fmt.Printf("Created a firewall rule called %s to %s, %s access to port %s %s %s with ID %s\n", utility.Green(rule.Label), utility.Green(newRuleConfig.Action), utility.Green(newRuleConfig.Direction), utility.Green(newRuleConfig.StartPort), directionValue, utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				} else {
					fmt.Printf("Created a firewall rule called %s to %s, %s access to ports %s-%s %s %s with ID %s\n", utility.Green(rule.Label), utility.Green(newRuleConfig.Action), utility.Green(newRuleConfig.Direction), utility.Green(newRuleConfig.StartPort), utility.Green(newRuleConfig.EndPort), directionValue, utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
				}
			}
		}
	},
}
