package loadbalancer

import (
	"errors"

	"github.com/spf13/cobra"
)

// LoadBalancerCmd manages Civo load balancers
var LoadBalancerCmd = &cobra.Command{
	Use:     "loadbalancer",
	Aliases: []string{"loadbalancers", "lb"},
	Short:   "Details of Civo Load Balancer",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	lbBackendExample := `Define backend configurations for the load balancer.
Use this format:
--backend "ip=10.0.0.1|protocol=http|source-port=80|target-port=8080|health-check-port=8080"
Each field is separated by '|'`

	lbInstancePoolExample := `Define instance pool configurations for the load balancer.
Use this format:
--instance-pool "tags=web,db|names=frontend,backend|protocol=TCP|source-port=80|target-port=8080|health-check.port=8080|health-check.path=/health"
Each field is separated by '|'`

	lbUpdateBackendExample := `Update backend configurations for the load balancer.
Use this format:
--backends "ip=10.0.0.1;protocol=http;source-port=80;target-port=8080;health-check-port=8080"
Each field is separated by ';'.
If you need multiple backends, use the flag multiple times.`

	lbUpdateInstancePoolExample := `Update instance pool configurations for the load balancer.
Use this format:
--instance-pools "tags=web,db;names=frontend,backend;protocol=TCP;source-port=80;target-port=8080;health-check.port=8080;health-check.path=/health"
Each field is separated by ';'.
If you need multiple instance pools, use the flag multiple times.`

	LoadBalancerCmd.AddCommand(loadBalancerListCmd)
	LoadBalancerCmd.AddCommand(loadBalancerShowCmd)
	LoadBalancerCmd.AddCommand(loadBalancerRemoveCmd)
	LoadBalancerCmd.AddCommand(loadBalancerCreateCmd)
	LoadBalancerCmd.AddCommand(loadBalancerUpdateCmd)

	// Balancer create subcommand
	loadBalancerCreateCmd.Flags().StringVarP(&lbName, "name", "", "", "Name of the load balancer")
	loadBalancerCreateCmd.Flags().StringVarP(&lbNetwork, "network", "n", "default", "The network to create the loadbalancer")
	loadBalancerCreateCmd.Flags().StringVarP(&lbAlgorithm, "algorithm", "a", "", "<round_robin | least_connections> - LoadBalancing algorithm to distribute traffic")
	loadBalancerCreateCmd.Flags().StringArrayVarP(&lbBackends, "backends", "b", []string{}, "Specify a backend instance to associate with the load balancer. Takes ip, protocol(optional), source-port, target-port and health-check-port(optional) in the format --backend=ip:instance-ip,protocol:HTTP|TCP,source-port:80,target-port:31579,health-check-port:31580")
	loadBalancerCreateCmd.Flags().StringVarP(&lbExternalTrafficPolicy, "external-traffic-policy", "e", "", "optional, Specify the external traffic policy for the load balancer")
	loadBalancerCreateCmd.Flags().StringVarP(&lbSessionAffinity, "session-affinity", "s", "", "optional, Specify the session affinity for the load balancer")
	loadBalancerCreateCmd.Flags().IntVarP(&lbSessionAffinityConfigTimeout, "session-affinity-config-timeout", "t", 0, "optional, Specify the session affinity config timeout for the load balancer")
	loadBalancerCreateCmd.Flags().StringVarP(&lbExistingFirewall, "existing-firewall", "v", "", "optional, ID of existing firewall to use")
	loadBalancerCreateCmd.Flags().StringVarP(&lbCreateFirewall, "create-firewall", "c", "", "optional, semicolon-separated list of ports to open - leave blank for default (80;443) or you can use \"all\"")
	loadBalancerCreateCmd.Flags().StringArrayVar(&lbBackends, "backend", []string{}, lbBackendExample)
	loadBalancerCreateCmd.Flags().StringSliceVar(&lbInstancePools, "instance-pool", []string{}, lbInstancePoolExample)

	// Balancer update subcommand
	loadBalancerUpdateCmd.Flags().StringVarP(&lbNameUpdate, "name", "", "", "New name of the load balancer to update (optional)")
	loadBalancerUpdateCmd.Flags().StringVarP(&lbAlgorithmUpdate, "algorithm", "a", "", "<round_robin | least_connections> - LoadBalancing algorithm to distribute traffic (optional)")
	loadBalancerUpdateCmd.Flags().StringArrayVarP(&lbBackendsUpdate, "backends", "b", []string{}, lbUpdateBackendExample)
	loadBalancerUpdateCmd.Flags().StringArrayVarP(&lbInstancePoolsUpdate, "instance-pools", "i", []string{}, lbUpdateInstancePoolExample)
}
