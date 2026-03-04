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

var vpcLoadBalancerShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: "civo vpc loadbalancer show ID/NAME",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show VPC load balancer",
	Long: `Show a specified VPC load balancer.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* algorithm
	* public_ip
	* state
	* private_ip
	* firewall_id
	* cluster_id
	* external_traffic_policy
	* session_affinity`,
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

		lb, err := client.FindVPCLoadBalancer(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", lb.ID, "ID")
		ow.AppendDataWithLabel("name", lb.Name, "Name")
		ow.AppendDataWithLabel("algorithm", lb.Algorithm, "Algorithm")
		ow.AppendDataWithLabel("public_ip", lb.PublicIP, "Public IP")
		ow.AppendDataWithLabel("state", lb.State, "State")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("private_ip", lb.PrivateIP, "Private IP")
			ow.AppendDataWithLabel("firewall_id", lb.FirewallID, "Firewall ID")
			ow.AppendDataWithLabel("cluster_id", lb.ClusterID, "Cluster ID")
			ow.AppendDataWithLabel("external_traffic_policy", lb.ExternalTrafficPolicy, "External Traffic Policy")
			ow.AppendDataWithLabel("session_affinity", lb.SessionAffinity, "Session Affinity")
		}

		var backendList []string
		for _, backend := range lb.Backends {
			backendList = append(backendList, backend.IP)
		}
		ow.AppendData("Backends", strings.Join(backendList, ", "))

		switch common.OutputFormat {
		case "json":
			ow.ToJSON(lb, common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Println("VPC Load Balancer Details:")
			fmt.Printf("ID: %s\n", lb.ID)
			fmt.Printf("Name: %s\n", lb.Name)
			fmt.Printf("Algorithm: %s\n", lb.Algorithm)
			fmt.Printf("Public IP: %s\n", lb.PublicIP)
			fmt.Printf("Private IP: %s\n", lb.PrivateIP)
			fmt.Printf("State: %s\n", lb.State)
			fmt.Printf("Firewall ID: %s\n", lb.FirewallID)
			fmt.Printf("External Traffic Policy: %s\n", lb.ExternalTrafficPolicy)
			fmt.Printf("Session Affinity: %s\n", lb.SessionAffinity)
			if len(lb.Backends) > 0 {
				fmt.Println("\nBackends:")
				for _, backend := range lb.Backends {
					fmt.Printf("  - %s (port %d -> %d, %s)\n", backend.IP, backend.SourcePort, backend.TargetPort, backend.Protocol)
				}
			}
		}
	},
}
