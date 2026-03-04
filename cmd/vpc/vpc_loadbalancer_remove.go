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

var vpcLBResourceList []utility.Resource
var vpcLoadBalancerRemoveCmd = &cobra.Command{
	Use:     "remove [NAME]",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo vpc loadbalancer remove NAME",
	Short:   "Remove a VPC load balancer",
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
			lb, err := client.FindVPCLoadBalancer(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s VPC load balancer in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one VPC load balancer with that name in your account")
					os.Exit(1)
				}
			}
			vpcLBResourceList = append(vpcLBResourceList, utility.Resource{ID: lb.ID, Name: lb.Name})
		} else {
			for _, v := range args {
				lb, err := client.FindVPCLoadBalancer(v)
				if err == nil {
					vpcLBResourceList = append(vpcLBResourceList, utility.Resource{ID: lb.ID, Name: lb.Name})
				}
			}
		}

		lbNameList := []string{}
		for _, v := range vpcLBResourceList {
			lbNameList = append(lbNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(vpcLBResourceList), "VPC load balancer"), common.DefaultYes, strings.Join(lbNameList, ", ")) {
			for _, v := range vpcLBResourceList {
				_, err = client.DeleteVPCLoadBalancer(v.ID)
				if err != nil {
					utility.Error("error deleting the VPC load balancer: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()
			for _, v := range vpcLBResourceList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(vpcLBResourceList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) %s been deleted\n",
					pluralize.Pluralize(len(vpcLBResourceList), "VPC load balancer"),
					utility.Green(strings.Join(lbNameList, ", ")),
					pluralize.Has(len(vpcLBResourceList)),
				)
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
