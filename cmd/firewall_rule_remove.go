package cmd

import (
	"fmt"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var firewallRuleRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "destroy", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove firewall rule",
	Example: "civo firewall rule remove FIREWALL_NAME/FIREWALL_ID FIREWALL_RULE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion("firewall", defaultRemove) == true {
			firewall, err := client.FindFirewall(args[0])
			if err != nil {
				utility.Error("Unable to find the firewall for your search %s", err)
				os.Exit(1)
			}

			rule, err := client.FindFirewallRule(firewall.ID, args[1])
			if err != nil {
				utility.Error("Unable to find the firewall rule %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteFirewallRule(firewall.ID, rule.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": rule.ID, "Label": rule.Label})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The firewall rule %s with ID %s was deleted\n", utility.Green(rule.Label), utility.Green(rule.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
