package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var firewallnetwork string
var defaultNetwork *civogo.Network

var firewallCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new firewall",
	Example: "civo firewall create NAME",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if firewallnetwork == "default" {
			defaultNetwork, err = client.GetDefaultNetwork()
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		} else {
			defaultNetwork, err = client.FindNetwork(firewallnetwork)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		firewall, err := client.NewFirewall(args[0], defaultNetwork.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": firewall.ID, "Name": firewall.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a firewall called %s with ID %s\n", utility.Green(firewall.Name), utility.Green(firewall.ID))
		}
	},
}
