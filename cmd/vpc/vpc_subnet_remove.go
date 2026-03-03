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

var subnetRemoveNetworkID string
var vpcSubnetResourceList []utility.Resource

var vpcSubnetRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo vpc subnet remove SUBNET_NAME --network NETWORK_NAME",
	Short:   "Remove a VPC subnet",
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

		network, err := client.FindVPCNetwork(subnetRemoveNetworkID)
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			subnet, err := client.FindVPCSubnet(args[0], network.ID)
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s VPC subnet in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one VPC subnet with that name in your account")
					os.Exit(1)
				}
			}
			vpcSubnetResourceList = append(vpcSubnetResourceList, utility.Resource{ID: subnet.ID, Name: subnet.Name})
		} else {
			for _, v := range args {
				subnet, err := client.FindVPCSubnet(v, network.ID)
				if err == nil {
					vpcSubnetResourceList = append(vpcSubnetResourceList, utility.Resource{ID: subnet.ID, Name: subnet.Name})
				}
			}
		}

		subnetNameList := []string{}
		for _, v := range vpcSubnetResourceList {
			subnetNameList = append(subnetNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(vpcSubnetResourceList), "VPC subnet"), common.DefaultYes, strings.Join(subnetNameList, ", ")) {
			for _, v := range vpcSubnetResourceList {
				_, err = client.DeleteVPCSubnet(network.ID, v.ID)
				if err != nil {
					utility.Error("Error deleting the VPC subnet: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()
			for _, v := range vpcSubnetResourceList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(vpcSubnetResourceList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) %s been deleted\n",
					pluralize.Pluralize(len(vpcSubnetResourceList), "VPC subnet"),
					utility.Green(strings.Join(subnetNameList, ", ")),
					pluralize.Has(len(vpcSubnetResourceList)),
				)
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
