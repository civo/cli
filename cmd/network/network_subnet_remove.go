package network

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

var subnetsList []utility.ObjecteList
var networkSubnetRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Short:   "Delete a Subnet from a network",
	Example: "civo network subnet delete <SUBNET-ID> <NETWORK-ID>",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		if len(args) == 2 {
			subnet, err := client.FindSubnet(args[0], args[1])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s subnet in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one subnet with that name in your account")
					os.Exit(1)
				}
			}
			subnetsList = append(subnetsList, utility.ObjecteList{ID: subnet.ID, Name: subnet.Name})
		} else {
			for _, v := range args {
				subnet, err := client.FindSubnet(v, args[1])
				if err == nil {
					subnetsList = append(subnetsList, utility.ObjecteList{ID: subnet.ID, Name: subnet.Name})
				}
			}
		}

		subnetNameList := []string{}
		for _, v := range subnetsList {
			subnetNameList = append(subnetNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(subnetNameList), "Subnet"), common.DefaultYes, strings.Join(subnetNameList, ", ")) {

			for _, v := range subnetsList {
				subnet, err := client.FindSubnet(v.ID, args[1])
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				_, err = client.DeleteSubnet(args[1], subnet.ID)
				if err != nil {
					utility.Error("Error deleting the Subnet: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range subnetsList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("subnet", v.Name, "Subnet")
			}

			switch common.OutputFormat {
			case "json":
				if len(subnetsList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(subnetsList), "subnet"), utility.Green(strings.Join(subnetNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
