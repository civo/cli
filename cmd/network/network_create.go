package network

import (
	"fmt"
	"os"

	"github.com/civo/civogo"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	cidrV4                string
	nameserversV4         []string
	createDefaultFirewall bool
)

var networkCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo network create NAME [flags]",
	Short:   "Create a new network",
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

		networkConfig := civogo.NetworkConfig{
			Label:         args[0],
			CIDRv4:        cidrV4,
			NameserversV4: nameserversV4,
			Region:        client.Region,
		}

		network, err := client.CreateNetwork(networkConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": network.ID, "label": network.Label})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Created a network called %s with ID %s\n", utility.Green(network.Label), utility.Green(network.ID))
		}

		if createDefaultFirewall {
			firewall, err := client.NewFirewall(&civogo.FirewallConfig{
				Name:      fmt.Sprintf("%s - Default", network.Label),
				NetworkID: network.ID,
				// CreateRules: &createRules,
				Region: client.Region,
			})
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			fmt.Printf("Created a default firewall %s\n", utility.Green(firewall.Name))
		}
	},
}
