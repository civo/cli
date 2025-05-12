package firewall

import (
	"errors"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var firewallList []utility.ObjecteList
var firewallRemoveCmd = &cobra.Command{
	Use:     "remove [NAME]",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo firewall remove NAME",
	Short:   "Remove a firewall",
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
			firewallList = append(firewallList, utility.ObjecteList{ID: firewall.ID, Name: firewall.Name})
		} else {
			for _, v := range args {
				firewall, err := client.FindFirewall(v)
				if err == nil {
					firewallList = append(firewallList, utility.ObjecteList{ID: firewall.ID, Name: firewall.Name})
				}
			}
		}

		firewallNameList := []string{}
		for _, v := range firewallList {
			firewallNameList = append(firewallNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(firewallList), "firewall"), common.DefaultYes, strings.Join(firewallNameList, ", ")) {

			for _, v := range firewallList {
				_, err = client.DeleteFirewall(v.ID)
				if err != nil {
					utility.Error("error deleting the firewall: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range firewallList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(firewallList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				utility.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(firewallList), "firewall"), utility.Green(strings.Join(firewallNameList, ", ")))
			}
		} else {
			utility.Println("Operation aborted.")
		}
	},
}
