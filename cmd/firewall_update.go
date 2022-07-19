package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var firewallUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"rename", "change"},
	Short:   "Update a firewall",
	Example: "civo firewall update OLD_NAME NEW_NAME",
	Args:    cobra.MinimumNArgs(2),
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

		defautlNetwork, err := client.GetDefaultNetwork()
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			utility.Error("Firewall %s", err)
			os.Exit(1)
		}

		firewallConfig := &civogo.FirewallConfig{
			Name:      args[1],
			NetworkID: defautlNetwork.ID,
		}

		_, err = client.RenameFirewall(firewall.ID, firewallConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": firewall.ID, "name": firewall.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The firewall called %s with ID %s was renamed to %s\n", utility.Green(firewall.Name), utility.Green(firewall.ID), utility.Green(args[1]))
		}
	},
}
