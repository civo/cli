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

var vpcFirewallResourceList []utility.Resource
var vpcFirewallRemoveCmd = &cobra.Command{
	Use:     "remove [NAME]",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo vpc firewall remove NAME",
	Short:   "Remove a VPC firewall",
	Args:    cobra.MinimumNArgs(1),
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

		if len(args) == 1 {
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
			vpcFirewallResourceList = append(vpcFirewallResourceList, utility.Resource{ID: firewall.ID, Name: firewall.Name})
		} else {
			for _, v := range args {
				firewall, err := client.FindVPCFirewall(v)
				if err == nil {
					vpcFirewallResourceList = append(vpcFirewallResourceList, utility.Resource{ID: firewall.ID, Name: firewall.Name})
				}
			}
		}

		firewallNameList := []string{}
		for _, v := range vpcFirewallResourceList {
			firewallNameList = append(firewallNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(vpcFirewallResourceList), "VPC firewall"), common.DefaultYes, strings.Join(firewallNameList, ", ")) {
			for _, v := range vpcFirewallResourceList {
				_, err = client.DeleteVPCFirewall(v.ID)
				if err != nil {
					utility.Error("error deleting the VPC firewall: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()
			for _, v := range vpcFirewallResourceList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(vpcFirewallResourceList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) %s been deleted\n",
					pluralize.Pluralize(len(vpcFirewallResourceList), "VPC firewall"),
					utility.Green(strings.Join(firewallNameList, ", ")),
					pluralize.Has(len(vpcFirewallResourceList)),
				)
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
