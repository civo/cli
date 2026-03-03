package vpc

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	vpcLBNetwork                string
	vpcLBAlgorithm              string
	vpcLBExternalTrafficPolicy  string
	vpcLBSessionAffinity        string
	vpcLBSessionAffinityTimeout int32
	vpcLBEnableProxyProtocol    string
	vpcLBFirewallID             string
	vpcLBMaxConcurrentRequests  int
	vpcLBBackends               []string
)

var vpcLoadBalancerCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: `civo vpc loadbalancer create NAME --network NETWORK [flags]
civo vpc loadbalancer create my-lb --network my-network --backend "10.0.0.1:80:8080:TCP:8080"`,
	Short: "Create a new VPC load balancer",
	Args:  cobra.MinimumNArgs(1),
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

		network, err := client.FindVPCNetwork(vpcLBNetwork)
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		lbConfig := &civogo.LoadBalancerConfig{
			Region:                client.Region,
			Name:                  args[0],
			NetworkID:             network.ID,
			Algorithm:             vpcLBAlgorithm,
			ExternalTrafficPolicy: vpcLBExternalTrafficPolicy,
			SessionAffinity:       vpcLBSessionAffinity,
			EnableProxyProtocol:   vpcLBEnableProxyProtocol,
			FirewallID:            vpcLBFirewallID,
		}

		if cmd.Flags().Changed("session-affinity-timeout") {
			lbConfig.SessionAffinityConfigTimeout = vpcLBSessionAffinityTimeout
		}

		if cmd.Flags().Changed("max-concurrent-requests") {
			lbConfig.MaxConcurrentRequests = &vpcLBMaxConcurrentRequests
		}

		if len(vpcLBBackends) > 0 {
			backends, err := parseBackends(vpcLBBackends)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			lbConfig.Backends = backends
		}

		lb, err := client.CreateVPCLoadBalancer(lbConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": lb.ID, "name": lb.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Created a VPC load balancer called %s with ID %s\n", utility.Green(lb.Name), utility.Green(lb.ID))
		}
	},
}

// parseBackends parses backend strings in format "ip:source_port:target_port:protocol:health_check_port"
func parseBackends(backends []string) ([]civogo.LoadBalancerBackendConfig, error) {
	var result []civogo.LoadBalancerBackendConfig
	for _, b := range backends {
		parts := strings.Split(b, ":")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid backend format '%s', expected ip:source_port:target_port[:protocol[:health_check_port]]", b)
		}

		sourcePort, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid source port '%s': %s", parts[1], err)
		}

		targetPort, err := strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid target port '%s': %s", parts[2], err)
		}

		backend := civogo.LoadBalancerBackendConfig{
			IP:         parts[0],
			SourcePort: int32(sourcePort),
			TargetPort: int32(targetPort),
		}

		if len(parts) > 3 {
			backend.Protocol = parts[3]
		}
		if len(parts) > 4 {
			healthPort, err := strconv.ParseInt(parts[4], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid health check port '%s': %s", parts[4], err)
			}
			backend.HealthCheckPort = int32(healthPort)
		}

		result = append(result, backend)
	}
	return result, nil
}
