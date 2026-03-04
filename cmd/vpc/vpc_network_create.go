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
	vpcNetCIDRv4        string
	vpcNetNameserversV4 []string
	vpcNetIPv4Enabled   bool
	vpcNetIPv6Enabled   bool
	vpcNetNameserversV6 []string
)

var vpcNetworkCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo vpc network create NAME [flags]",
	Short:   "Create a new VPC network",
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
			Label:  args[0],
			Region: client.Region,
		}

		if vpcNetCIDRv4 != "" {
			networkConfig.CIDRv4 = vpcNetCIDRv4
		}
		if len(vpcNetNameserversV4) > 0 {
			networkConfig.NameserversV4 = vpcNetNameserversV4
		}
		if cmd.Flags().Changed("ipv4-enabled") {
			networkConfig.IPv4Enabled = &vpcNetIPv4Enabled
		}
		if cmd.Flags().Changed("ipv6-enabled") {
			networkConfig.IPv6Enabled = &vpcNetIPv6Enabled
		}
		if len(vpcNetNameserversV6) > 0 {
			networkConfig.NameserversV6 = vpcNetNameserversV6
		}

		network, err := client.CreateVPCNetwork(networkConfig)
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
			fmt.Printf("Created a VPC network called %s with ID %s\n", utility.Green(network.Label), utility.Green(network.ID))
		}
	},
}
