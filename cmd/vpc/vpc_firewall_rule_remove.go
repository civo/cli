package vpc

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/pkg/pluralize"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcFirewallRuleResourceList []utility.Resource
var vpcFirewallRuleRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "destroy", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove VPC firewall rule",
	Example: "civo vpc firewall rule remove FIREWALL_NAME/FIREWALL_ID FIREWALL_RULE_ID",
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

		firewall, err := client.FindVPCFirewall(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s VPC firewall in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one VPC firewall with that name in your account")
				os.Exit(1)
			}
		}

		if len(args) == 2 {
			rule, err := client.FindVPCFirewallRule(firewall.ID, args[1])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s VPC firewall rule in your account", utility.Red(args[1]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one VPC firewall rule in your account")
					os.Exit(1)
				}
			}
			vpcFirewallRuleResourceList = append(vpcFirewallRuleResourceList, utility.Resource{ID: rule.ID, Name: rule.Label})
		} else {
			for _, v := range args[1:] {
				rule, err := client.FindVPCFirewallRule(firewall.ID, v)
				if err == nil {
					vpcFirewallRuleResourceList = append(vpcFirewallRuleResourceList, utility.Resource{ID: rule.ID, Name: rule.Label})
				}
			}
		}

		firewallRuleNameList := []string{}
		for _, v := range vpcFirewallRuleResourceList {
			firewallRuleNameList = append(firewallRuleNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("VPC firewall %s", pluralize.Pluralize(len(vpcFirewallRuleResourceList), "rule")), common.DefaultYes, strings.Join(firewallRuleNameList, ", ")) {
			for _, v := range vpcFirewallRuleResourceList {
				_, err = client.DeleteVPCFirewallRule(firewall.ID, v.ID)
				if err != nil {
					utility.Error("error deleting the VPC firewall rule: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()
			for _, v := range vpcFirewallRuleResourceList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("label", v.Name, "Label")
			}

			switch common.OutputFormat {
			case "json":
				if len(vpcFirewallRuleResourceList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The VPC firewall %s (%s) %s been deleted\n",
					pluralize.Pluralize(len(vpcFirewallRuleResourceList), "rule"),
					strings.Join(firewallRuleNameList, ", "),
					pluralize.Has(len(vpcFirewallRuleResourceList)),
				)
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
