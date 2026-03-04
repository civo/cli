package vpc

import (
	"errors"

	"github.com/spf13/cobra"
)

// VPCCmd manages Civo VPC resources
var VPCCmd = &cobra.Command{
	Use:     "vpc",
	Aliases: []string{},
	Short:   "Details of Civo VPC resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	// Sub-resource commands
	VPCCmd.AddCommand(vpcNetworkCmd)
	VPCCmd.AddCommand(vpcSubnetCmd)
	VPCCmd.AddCommand(vpcFirewallCmd)
	VPCCmd.AddCommand(vpcLoadBalancerCmd)
	VPCCmd.AddCommand(vpcIPCmd)

	// --- Network subcommands ---
	vpcNetworkCmd.AddCommand(vpcNetworkListCmd)
	vpcNetworkCmd.AddCommand(vpcNetworkCreateCmd)
	vpcNetworkCmd.AddCommand(vpcNetworkShowCmd)
	vpcNetworkCmd.AddCommand(vpcNetworkUpdateCmd)
	vpcNetworkCmd.AddCommand(vpcNetworkRemoveCmd)

	// Network create flags
	vpcNetworkCreateCmd.Flags().StringVarP(&vpcNetCIDRv4, "cidr-v4", "", "", "Custom IPv4 CIDR")
	vpcNetworkCreateCmd.Flags().StringSliceVarP(&vpcNetNameserversV4, "nameservers-v4", "", nil, "Custom list of IPv4 nameservers (comma-separated)")
	vpcNetworkCreateCmd.Flags().BoolVarP(&vpcNetIPv4Enabled, "ipv4-enabled", "", true, "Enable IPv4 for the network")
	vpcNetworkCreateCmd.Flags().BoolVarP(&vpcNetIPv6Enabled, "ipv6-enabled", "", false, "Enable IPv6 for the network")
	vpcNetworkCreateCmd.Flags().StringSliceVarP(&vpcNetNameserversV6, "nameservers-v6", "", nil, "Custom list of IPv6 nameservers (comma-separated)")

	// --- Subnet subcommands ---
	vpcSubnetCmd.AddCommand(vpcSubnetListCmd)
	vpcSubnetCmd.AddCommand(vpcSubnetCreateCmd)
	vpcSubnetCmd.AddCommand(vpcSubnetShowCmd)
	vpcSubnetCmd.AddCommand(vpcSubnetRemoveCmd)
	vpcSubnetCmd.AddCommand(vpcSubnetAttachCmd)
	vpcSubnetCmd.AddCommand(vpcSubnetDetachCmd)

	// Subnet flags (all require --network)
	vpcSubnetListCmd.Flags().StringVarP(&subnetListNetworkID, "network", "n", "", "Network name or ID")
	_ = vpcSubnetListCmd.MarkFlagRequired("network")

	vpcSubnetCreateCmd.Flags().StringVarP(&subnetCreateNetworkID, "network", "n", "", "Network name or ID")
	_ = vpcSubnetCreateCmd.MarkFlagRequired("network")

	vpcSubnetShowCmd.Flags().StringVarP(&subnetShowNetworkID, "network", "n", "", "Network name or ID")
	_ = vpcSubnetShowCmd.MarkFlagRequired("network")

	vpcSubnetRemoveCmd.Flags().StringVarP(&subnetRemoveNetworkID, "network", "n", "", "Network name or ID")
	_ = vpcSubnetRemoveCmd.MarkFlagRequired("network")

	vpcSubnetAttachCmd.Flags().StringVarP(&subnetAttachNetworkID, "network", "n", "", "Network name or ID")
	_ = vpcSubnetAttachCmd.MarkFlagRequired("network")
	vpcSubnetAttachCmd.Flags().StringVar(&subnetAttachResourceID, "resource-id", "", "Resource ID to attach to")
	_ = vpcSubnetAttachCmd.MarkFlagRequired("resource-id")
	vpcSubnetAttachCmd.Flags().StringVar(&subnetAttachResourceType, "resource-type", "", "Resource type (e.g. instance)")
	_ = vpcSubnetAttachCmd.MarkFlagRequired("resource-type")

	vpcSubnetDetachCmd.Flags().StringVarP(&subnetDetachNetworkID, "network", "n", "", "Network name or ID")
	_ = vpcSubnetDetachCmd.MarkFlagRequired("network")

	// --- Firewall subcommands ---
	vpcFirewallCmd.AddCommand(vpcFirewallListCmd)
	vpcFirewallCmd.AddCommand(vpcFirewallCreateCmd)
	vpcFirewallCmd.AddCommand(vpcFirewallShowCmd)
	vpcFirewallCmd.AddCommand(vpcFirewallUpdateCmd)
	vpcFirewallCmd.AddCommand(vpcFirewallRemoveCmd)

	// Firewall create flags
	vpcFirewallCreateCmd.Flags().StringVarP(&vpcFwNetwork, "network", "n", "default", "The network to create the firewall in")
	vpcFirewallCreateCmd.Flags().BoolVarP(&vpcFwNoDefaultRules, "no-default-rules", "", false, "Create firewall without default rules")

	// --- Firewall Rule subcommands ---
	vpcFirewallCmd.AddCommand(vpcFirewallRuleCmd)
	vpcFirewallRuleCmd.AddCommand(vpcFirewallRuleListCmd)
	vpcFirewallRuleCmd.AddCommand(vpcFirewallRuleCreateCmd)
	vpcFirewallRuleCmd.AddCommand(vpcFirewallRuleRemoveCmd)

	// Firewall rule create flags
	vpcFirewallRuleCreateCmd.Flags().StringVarP(&vpcFwRuleProtocol, "protocol", "p", "TCP", "The protocol choice (TCP, UDP, ICMP)")
	vpcFirewallRuleCreateCmd.Flags().StringVarP(&vpcFwRuleStartPort, "startport", "s", "", "The start port of the rule")
	vpcFirewallRuleCreateCmd.Flags().StringVarP(&vpcFwRuleEndPort, "endport", "e", "", "The end port of the rule")
	vpcFirewallRuleCreateCmd.Flags().StringVarP(&vpcFwRuleCidr, "cidr", "c", "0.0.0.0/0", "The CIDR of the rule")
	vpcFirewallRuleCreateCmd.Flags().StringVarP(&vpcFwRuleDirection, "direction", "d", "ingress", "The direction of the rule (ingress or egress)")
	vpcFirewallRuleCreateCmd.Flags().StringVarP(&vpcFwRuleAction, "action", "a", "allow", "The action of the rule (allow or deny)")
	vpcFirewallRuleCreateCmd.Flags().StringVarP(&vpcFwRuleLabel, "label", "l", "", "A label for this rule")
	_ = vpcFirewallRuleCreateCmd.MarkFlagRequired("startport")

	// --- Load Balancer subcommands ---
	vpcLoadBalancerCmd.AddCommand(vpcLoadBalancerListCmd)
	vpcLoadBalancerCmd.AddCommand(vpcLoadBalancerCreateCmd)
	vpcLoadBalancerCmd.AddCommand(vpcLoadBalancerShowCmd)
	vpcLoadBalancerCmd.AddCommand(vpcLoadBalancerUpdateCmd)
	vpcLoadBalancerCmd.AddCommand(vpcLoadBalancerRemoveCmd)

	// Load balancer create flags
	vpcLoadBalancerCreateCmd.Flags().StringVarP(&vpcLBNetwork, "network", "n", "", "Network name or ID")
	_ = vpcLoadBalancerCreateCmd.MarkFlagRequired("network")
	vpcLoadBalancerCreateCmd.Flags().StringVar(&vpcLBAlgorithm, "algorithm", "round_robin", "Load balancing algorithm")
	vpcLoadBalancerCreateCmd.Flags().StringVar(&vpcLBExternalTrafficPolicy, "external-traffic-policy", "", "External traffic policy")
	vpcLoadBalancerCreateCmd.Flags().StringVar(&vpcLBSessionAffinity, "session-affinity", "", "Session affinity")
	vpcLoadBalancerCreateCmd.Flags().Int32Var(&vpcLBSessionAffinityTimeout, "session-affinity-timeout", 0, "Session affinity timeout in seconds")
	vpcLoadBalancerCreateCmd.Flags().StringVar(&vpcLBEnableProxyProtocol, "enable-proxy-protocol", "", "Enable proxy protocol")
	vpcLoadBalancerCreateCmd.Flags().StringVar(&vpcLBFirewallID, "firewall-id", "", "Firewall ID to associate")
	vpcLoadBalancerCreateCmd.Flags().IntVar(&vpcLBMaxConcurrentRequests, "max-concurrent-requests", 0, "Maximum concurrent requests")
	vpcLoadBalancerCreateCmd.Flags().StringSliceVar(&vpcLBBackends, "backend", nil, "Backend in format ip:source_port:target_port:protocol:health_check_port (repeatable)")

	// Load balancer update flags
	vpcLoadBalancerUpdateCmd.Flags().StringVar(&vpcLBUpdateName, "name", "", "New name for the load balancer")
	vpcLoadBalancerUpdateCmd.Flags().StringVar(&vpcLBUpdateAlgorithm, "algorithm", "", "Load balancing algorithm")
	vpcLoadBalancerUpdateCmd.Flags().StringVar(&vpcLBUpdateExternalTrafficPolicy, "external-traffic-policy", "", "External traffic policy")
	vpcLoadBalancerUpdateCmd.Flags().StringVar(&vpcLBUpdateSessionAffinity, "session-affinity", "", "Session affinity")
	vpcLoadBalancerUpdateCmd.Flags().Int32Var(&vpcLBUpdateSessionAffinityTimeout, "session-affinity-timeout", 0, "Session affinity timeout in seconds")
	vpcLoadBalancerUpdateCmd.Flags().StringVar(&vpcLBUpdateEnableProxyProtocol, "enable-proxy-protocol", "", "Enable proxy protocol")
	vpcLoadBalancerUpdateCmd.Flags().StringVar(&vpcLBUpdateFirewallID, "firewall-id", "", "Firewall ID to associate")
	vpcLoadBalancerUpdateCmd.Flags().IntVar(&vpcLBUpdateMaxConcurrentRequests, "max-concurrent-requests", 0, "Maximum concurrent requests")
	vpcLoadBalancerUpdateCmd.Flags().StringSliceVar(&vpcLBUpdateBackends, "backend", nil, "Backend in format ip:source_port:target_port:protocol:health_check_port (repeatable)")

	// --- IP subcommands ---
	vpcIPCmd.AddCommand(vpcIPListCmd)
	vpcIPCmd.AddCommand(vpcIPCreateCmd)
	vpcIPCmd.AddCommand(vpcIPShowCmd)
	vpcIPCmd.AddCommand(vpcIPUpdateCmd)
	vpcIPCmd.AddCommand(vpcIPAssignCmd)
	vpcIPCmd.AddCommand(vpcIPUnassignCmd)
	vpcIPCmd.AddCommand(vpcIPRemoveCmd)

	// IP create flags
	vpcIPCreateCmd.Flags().StringVarP(&vpcIPName, "name", "n", "", "Name of the reserved IP")

	// IP update flags
	vpcIPUpdateCmd.Flags().StringVarP(&vpcIPUpdateName, "name", "n", "", "New name for the reserved IP")
	_ = vpcIPUpdateCmd.MarkFlagRequired("name")

	// IP assign flags
	vpcIPAssignCmd.Flags().StringVar(&vpcIPAssignResourceID, "resource-id", "", "Resource ID to assign the IP to")
	_ = vpcIPAssignCmd.MarkFlagRequired("resource-id")
	vpcIPAssignCmd.Flags().StringVar(&vpcIPAssignResourceType, "resource-type", "", "Resource type (e.g. instance, loadbalancer)")
	_ = vpcIPAssignCmd.MarkFlagRequired("resource-type")
}
