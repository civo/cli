package vpc

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcLoadBalancerListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo vpc loadbalancer ls -o custom -f "ID: Name"`,
	Short:   "List VPC load balancers",
	Long: `List all VPC load balancers.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* algorithm
	* public_ip
	* private_ip
	* state
	* backends`,
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

		lbs, err := client.ListVPCLoadBalancers()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, lb := range lbs {
			ow.StartLine()
			ow.AppendDataWithLabel("id", lb.ID, "ID")
			ow.AppendDataWithLabel("name", lb.Name, "Name")
			ow.AppendDataWithLabel("algorithm", lb.Algorithm, "Algorithm")
			ow.AppendDataWithLabel("public_ip", lb.PublicIP, "Public IP")
			ow.AppendDataWithLabel("private_ip", lb.PrivateIP, "Private IP")
			ow.AppendDataWithLabel("state", lb.State, "State")

			var backendList []string
			for _, backend := range lb.Backends {
				backendList = append(backendList, backend.IP)
			}
			ow.AppendDataWithLabel("backends", fmt.Sprintf("%d", len(lb.Backends)), "Backends")

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("firewall_id", lb.FirewallID, "Firewall ID")
				ow.AppendDataWithLabel("cluster_id", lb.ClusterID, "Cluster ID")
				ow.AppendDataWithLabel("external_traffic_policy", lb.ExternalTrafficPolicy, "External Traffic Policy")
				ow.AppendDataWithLabel("session_affinity", lb.SessionAffinity, "Session Affinity")
				ow.AppendData("Backend IPs", strings.Join(backendList, ", "))
			}
		}

		ow.FinishAndPrintOutput()
	},
}
