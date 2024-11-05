package loadbalancer

import (
	"fmt"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var loadBalancerShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo loadbalancer show ID/NAME -o custom -f "ID: Name"`,
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show load balancer",
	Long: `Show a specified load balancer.
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
	* session_affinity
	* session_affinity_config_timeout
	* dns_entry`,
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

		lb, err := client.GetLoadBalancer(args[0])
		if err != nil {
			utility.Error("Failed to retrieve load balancer: %s", err)
			os.Exit(1)
		}

		// Display Core Load Balancer Details
		fmt.Printf("ID: %s\n", lb.ID)
		fmt.Printf("Name: %s\n", lb.Name)
		fmt.Printf("Algorithm: %s\n", lb.Algorithm)
		fmt.Printf("Public IP: %s\n", lb.PublicIP)
		fmt.Printf("Private IP: %s\n", lb.PrivateIP)
		fmt.Printf("Firewall ID: %s\n", lb.FirewallID)
		fmt.Printf("Cluster ID: %s\n", lb.ClusterID)
		fmt.Printf("State: %s\n", lb.State)
		fmt.Printf("DNS Entry: %s.lb.civo.com\n", lb.ID)

		// Session and Traffic Policies
		fmt.Println("\nPolicies:")
		fmt.Printf("External Traffic Policy: %s\n", lb.ExternalTrafficPolicy)
		fmt.Printf("Session Affinity: %s\n", lb.SessionAffinity)
		fmt.Printf("Session Affinity Config Timeout: %d\n", lb.SessionAffinityConfigTimeout)
		fmt.Printf("Enable Proxy Protocol: %s\n", lb.EnableProxyProtocol)

		// Reserved IP Information
		if lb.ReservedIPID != "" || lb.ReservedIP != "" || lb.ReservedIPName != "" {
			fmt.Println("\nReserved IP Details:")
			fmt.Printf("Reserved IP ID: %s\n", lb.ReservedIPID)
			fmt.Printf("Reserved IP Name: %s\n", lb.ReservedIPName)
			fmt.Printf("Reserved IP: %s\n", lb.ReservedIP)
		} else {
			fmt.Println("\nNo Reserved IP Configuration")
		}

		// Options if available
		if lb.Options != nil {
			fmt.Println("\nLoad Balancer Options:")
			fmt.Printf("Server Timeout: %s\n", lb.Options.ServerTimeout)
			fmt.Printf("Client Timeout: %s\n", lb.Options.ClientTimeout)
		}

		// Display Backends
		if len(lb.Backends) > 0 {
			fmt.Println("\nBackends:")
			for _, backend := range lb.Backends {
				fmt.Printf("- IP: %s, Protocol: %s, Source Port: %d, Target Port: %d\n",
					backend.IP, backend.Protocol, backend.SourcePort, backend.TargetPort)
			}
		} else {
			fmt.Println("\nNo Backends Configured")
		}

		// Display Instance Pools
		if len(lb.InstancePool) > 0 {
			fmt.Println("\nInstance Pools:")
			for _, pool := range lb.InstancePool {
				fmt.Printf(" - Tags: %s, Protocol: %s, Source Port: %d, Target Port: %d, Health Check Port: %d, Health Check Path: %s\n",
					pool.Tags, pool.Protocol, pool.SourcePort, pool.TargetPort, pool.HealthCheck.Port, pool.HealthCheck.Path)
			}
		} else {
			fmt.Println("\nNo Instance Pools Configured")
		}

		// Max Concurrent Requests
		fmt.Printf("\nMax Concurrent Requests: %d\n", lb.MaxConcurrentRequests)
	},
}
