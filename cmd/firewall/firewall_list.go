package firewall

import (
	"os"
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var firewallListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List firewall",
	Long: `List all current firewalls.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* network
	* rules_count
	* instances_count

Example: civo firewall ls -o custom -f "ID: Name"`,
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

		firewalls, err := client.ListFirewalls()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		networks, err := client.ListNetworks()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		networkMap := make(map[string]string, 0)
		for _, network := range networks {
			networkMap[network.ID] = network.Label
		}

		ow := utility.NewOutputWriter()
		for _, firewall := range firewalls {

			ow.StartLine()

			ow.AppendDataWithLabel("id", firewall.ID, "ID")
			ow.AppendDataWithLabel("name", firewall.Name, "Name")
			ow.AppendDataWithLabel("network", networkMap[firewall.NetworkID], "Network")
			ow.AppendDataWithLabel("rules_count", strconv.Itoa(firewall.RulesCount), "Total rules")
			ow.AppendDataWithLabel("instances_count", strconv.Itoa(firewall.InstanceCount), "Total Instances")
			ow.AppendDataWithLabel("clusters_count", strconv.Itoa(firewall.ClusterCount), "Total Clusters")
			ow.AppendDataWithLabel("loadbalancer_count", strconv.Itoa(firewall.LoadBalancerCount), "Total LoadBalancer")
		}

		ow.FinishAndPrintOutput()
	},
}
