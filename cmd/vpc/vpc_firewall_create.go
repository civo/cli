package vpc

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
	vpcFwNetwork        string
	vpcFwNoDefaultRules bool
)

var vpcFirewallCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new VPC firewall",
	Example: "civo vpc firewall create NAME [flags]",
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

		var networkID string
		if vpcFwNetwork == "default" {
			defaultNetwork, err := client.GetDefaultVPCNetwork()
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}
			networkID = defaultNetwork.ID
		} else {
			network, err := client.FindVPCNetwork(vpcFwNetwork)
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}
			networkID = network.ID
		}

		createRules := !vpcFwNoDefaultRules

		firewall, err := client.NewVPCFirewall(&civogo.FirewallConfig{
			Name:        args[0],
			NetworkID:   networkID,
			CreateRules: &createRules,
			Region:      client.Region,
		})
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
			fmt.Printf("Created a VPC firewall called %s with ID %s\n", utility.Green(firewall.Name), utility.Green(firewall.ID))
		}
	},
}
