package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var firewallRuleRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "destroy", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove firewall rule",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("firewall") == nil {
			firewall, err := client.FindFirewall(args[0])
			if err != nil {
				fmt.Printf("Unable to find the firewall for your search: %s\n", aurora.Red(err))
				os.Exit(1)
			}

			rule, err := client.FindFirewallRule(firewall.ID, args[1])
			if err != nil {
				fmt.Printf("Unable to find the firewall rule: %s\n", aurora.Red(err))
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
				fmt.Printf("The firewall rule %s with ID %s was delete\n", aurora.Green(rule.Label), aurora.Green(rule.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
