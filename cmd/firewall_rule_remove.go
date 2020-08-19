package cmd

import (
	"errors"
	"fmt"

	"github.com/civo/civogo"
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

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry this firewall (%s) does not exist in your account", args[0])
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one firewall with that name in your account", args[0])
				os.Exit(1)
			}
		}

		rule, err := client.FindFirewallRule(firewall.ID, args[1])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry this firewall rule (%s) does not exist in your account", args[1])
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one firewall rule in your account", args[1])
				os.Exit(1)
			}
		}

		if utility.UserConfirmedDeletion("firewall", defaultYes) == true {

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
