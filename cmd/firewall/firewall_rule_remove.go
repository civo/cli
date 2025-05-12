package firewall

import (
	"errors"
	"fmt"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var firewallRuleList []utility.ObjecteList
var firewallRuleRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "destroy", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove firewall rule",
	Example: "civo firewall rule remove FIREWALL_NAME/FIREWALL_ID FIREWALL_RULE_ID",
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
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s firewall in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one firewall with that name in your account")
				os.Exit(1)
			}
		}

		if len(args) == 2 {
			rule, err := client.FindFirewallRule(firewall.ID, args[1])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s firewall rule in your account", utility.Red(args[1]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one firewall rule in your account")
					os.Exit(1)
				}
			}
			firewallRuleList = append(firewallRuleList, utility.ObjecteList{ID: rule.ID, Name: rule.Label})
		} else {
			for _, v := range args[1:] {
				rule, err := client.FindFirewallRule(firewall.ID, v)
				if err == nil {
					firewallRuleList = append(firewallRuleList, utility.ObjecteList{ID: rule.ID, Name: rule.Label})
				}
			}
		}

		firewallRuleNameList := []string{}
		for _, v := range firewallRuleList {
			firewallRuleNameList = append(firewallRuleNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("firewall %s", pluralize.Pluralize(len(firewallRuleList), "rule")), common.DefaultYes, strings.Join(firewallRuleNameList, ", ")) {

			for _, v := range firewallRuleList {
				_, err = client.DeleteFirewallRule(firewall.ID, v.ID)
				if err != nil {
					utility.Error("error deleting the firewall rule: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range firewallRuleList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("label", v.Name, "Label")
			}

			switch common.OutputFormat {
			case "json":
				if len(firewallRuleList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				utility.Printf("The firewall %s (%s) has been deleted\n", pluralize.Pluralize(len(firewallRuleList), "rule"), strings.Join(firewallRuleNameList, ", "))
			}
		} else {
			utility.Println("Operation aborted.")
		}

	},
}
