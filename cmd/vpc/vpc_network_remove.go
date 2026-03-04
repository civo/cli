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

var vpcNetworkResourceList []utility.Resource
var vpcNetworkRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo vpc network remove NAME",
	Short:   "Remove a VPC network",
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
			network, err := client.FindVPCNetwork(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s VPC network in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one VPC network with that name in your account")
					os.Exit(1)
				}
			}
			vpcNetworkResourceList = append(vpcNetworkResourceList, utility.Resource{ID: network.ID, Name: network.Label})
		} else {
			for _, v := range args {
				network, err := client.FindVPCNetwork(v)
				if err != nil {
					if errors.Is(err, civogo.ZeroMatchesError) {
						utility.Error("sorry there is no %s VPC network in your account", utility.Red(v))
						os.Exit(1)
					}
					if errors.Is(err, civogo.MultipleMatchesError) {
						utility.Error("sorry we found more than one VPC network with that name in your account")
						os.Exit(1)
					}
				}
				if err == nil {
					vpcNetworkResourceList = append(vpcNetworkResourceList, utility.Resource{ID: network.ID, Name: network.Label})
				}
			}
		}

		networkNameList := []string{}
		for _, v := range vpcNetworkResourceList {
			networkNameList = append(networkNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(vpcNetworkResourceList), "VPC network"), common.DefaultYes, strings.Join(networkNameList, ", ")) {
			for _, v := range vpcNetworkResourceList {
				_, err = client.DeleteVPCNetwork(v.ID)
				if err != nil {
					utility.Error("Error deleting the VPC network: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()
			for _, v := range vpcNetworkResourceList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("label", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(vpcNetworkResourceList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) %s been deleted\n",
					pluralize.Pluralize(len(vpcNetworkResourceList), "VPC network"),
					utility.Green(strings.Join(networkNameList, ", ")),
					pluralize.Has(len(vpcNetworkResourceList)),
				)
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
