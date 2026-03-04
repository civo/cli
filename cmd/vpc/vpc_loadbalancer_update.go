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
	vpcLBUpdateName                   string
	vpcLBUpdateAlgorithm              string
	vpcLBUpdateExternalTrafficPolicy  string
	vpcLBUpdateSessionAffinity        string
	vpcLBUpdateSessionAffinityTimeout int32
	vpcLBUpdateEnableProxyProtocol    string
	vpcLBUpdateFirewallID             string
	vpcLBUpdateMaxConcurrentRequests  int
	vpcLBUpdateBackends               []string
)

var vpcLoadBalancerUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Example: "civo vpc loadbalancer update LB_NAME/LB_ID [flags]",
	Short:   "Update a VPC load balancer",
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

		lb, err := client.FindVPCLoadBalancer(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		updateConfig := &civogo.LoadBalancerUpdateConfig{
			Region: client.Region,
		}

		if cmd.Flags().Changed("name") {
			updateConfig.Name = vpcLBUpdateName
		}
		if cmd.Flags().Changed("algorithm") {
			updateConfig.Algorithm = vpcLBUpdateAlgorithm
		}
		if cmd.Flags().Changed("external-traffic-policy") {
			updateConfig.ExternalTrafficPolicy = vpcLBUpdateExternalTrafficPolicy
		}
		if cmd.Flags().Changed("session-affinity") {
			updateConfig.SessionAffinity = vpcLBUpdateSessionAffinity
		}
		if cmd.Flags().Changed("session-affinity-timeout") {
			updateConfig.SessionAffinityConfigTimeout = vpcLBUpdateSessionAffinityTimeout
		}
		if cmd.Flags().Changed("enable-proxy-protocol") {
			updateConfig.EnableProxyProtocol = vpcLBUpdateEnableProxyProtocol
		}
		if cmd.Flags().Changed("firewall-id") {
			updateConfig.FirewallID = vpcLBUpdateFirewallID
		}
		if cmd.Flags().Changed("max-concurrent-requests") {
			updateConfig.MaxConcurrentRequests = &vpcLBUpdateMaxConcurrentRequests
		}
		if len(vpcLBUpdateBackends) > 0 {
			backends, err := parseBackends(vpcLBUpdateBackends)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			updateConfig.Backends = backends
		}

		updatedLB, err := client.UpdateVPCLoadBalancer(lb.ID, updateConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": updatedLB.ID, "name": updatedLB.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Updated VPC load balancer %s with ID %s\n", utility.Green(updatedLB.Name), utility.Green(updatedLB.ID))
		}
	},
}
