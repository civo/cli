package network

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var networkList []utility.ObjecteList
var networkRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo network rm NAME",
	Short:   "Remove a network",
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
			network, err := client.FindNetwork(args[0])

			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s network in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one network with that name in your account")
					os.Exit(1)
				}
			}

			networkList = append(networkList, utility.ObjecteList{ID: network.ID, Name: network.Label})

		} else {
			for _, v := range args {
				network, err := client.FindNetwork(v)
				if err != nil {
					if errors.Is(err, civogo.ZeroMatchesError) {
						utility.Error("sorry there is no %s network in your account", utility.Red(args[0]))
						os.Exit(1)
					}
					if errors.Is(err, civogo.MultipleMatchesError) {
						utility.Error("sorry we found more than one network with that name in your account")
						os.Exit(1)
					}
				}
				if err == nil {
					networkList = append(networkList, utility.ObjecteList{ID: network.ID, Name: network.Label})
				}
			}
		}

		networkNameList := []string{}
		for _, v := range networkList {
			networkNameList = append(networkNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(networkList), "network"), common.DefaultYes, strings.Join(networkNameList, ", ")) {

			for _, v := range networkList {
				_, err = client.DeleteNetwork(v.ID)
				if err != nil {
					if errors.Is(err, civogo.DatabaseNetworkDeleteWithInstanceError) {
						errMessage := fmt.Sprintf("sorry couldn't delete this network %s while it is in use\n", utility.Green(v.Name))
						utility.Error(errMessage)
						os.Exit(1)
					}
					utility.Error("Error deleting the Network: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range networkList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("label", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(networkList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(networkList), "network"), utility.Green(strings.Join(networkNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
