package loadbalancer

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

var lbName, lbNetwork, lbAlgorithm, lbExternalTrafficPolicy, lbSessionAffinity, lbExistingFirewall, lbCreateFirewall string
var lbSessionAffinityConfigTimeout int
var lbBackends []string
var lbInstancePools []string

var loadBalancerCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: `civo loadbalancer create my-loadbalancer \
    --network default \
    --create-firewall "80;443" \
    --algorithm "round-robin" \
    --session-affinity "ClientIP" \
    --session-affinity-config-timeout 10800 \
    --external-traffic-policy "Local" \
    --backend "ip:10.0.0.1,source-port:80,target-port:8080,protocol:http,health-check-port:8080,protocol:TCP" \
    --backend "ip:10.0.0.2,source-port:80,target-port:8080,protocol:http,health-check-port:8080,protocol:TCP"`,
	Short: "Create a new load balancer",
	Run: func(cmd *cobra.Command, args []string) {
		runLoadBalancerCreate(args)
	},
}

func runLoadBalancerCreate(args []string) {
	utility.EnsureCurrentRegion()

	check, region, err := utility.CheckAvailability("iaas", common.RegionSet)
	handleAvailabilityCheck(check, region, err)

	client := getCivoClient()
	configLoadBalancer := &civogo.LoadBalancerConfig{}

	setLoadBalancerName(configLoadBalancer, args)
	setLoadBalancerNetwork(client, configLoadBalancer)
	setLoadBalancerOptions(configLoadBalancer)

	if len(lbBackends) > 0 {
		err := setLoadBalancerBackends(configLoadBalancer)
		if err != nil {
			utility.Error(err.Error())
			os.Exit(1)
		}
	}

	if len(lbInstancePools) > 0 {
		err := setLoadBalancerInstancePools(configLoadBalancer)
		if err != nil {
			utility.Error(err.Error())
			os.Exit(1)
		}
	}

	loadBalancer, err := client.CreateLoadBalancer(configLoadBalancer)
	if err != nil {
		utility.Error("Creating the load balancer failed with %s", err)
		os.Exit(1)
	}

	outputLoadBalancer(loadBalancer)
}

func handleAvailabilityCheck(check bool, region string, err error) {
	if err != nil {
		utility.Error("Error checking availability %s", err)
		os.Exit(1)
	}

	if check {
		utility.Error("Sorry you can't create a load balancer in the %s region", region)
		os.Exit(1)
	}
}

func getCivoClient() *civogo.Client {
	client, err := config.CivoAPIClient()
	if common.RegionSet != "" {
		client.Region = common.RegionSet
	}
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}
	return client
}

func setLoadBalancerName(configLoadBalancer *civogo.LoadBalancerConfig, args []string) {
	if len(args) > 0 {
		if utility.ValidNameLength(args[0]) {
			utility.Warning("The load balancer name cannot be longer than 63 characters")
			os.Exit(1)
		}
		configLoadBalancer.Name = args[0]
	} else {
		configLoadBalancer.Name = utility.RandomName()
	}
}

func setLoadBalancerNetwork(client *civogo.Client, configLoadBalancer *civogo.LoadBalancerConfig) {
	var network *civogo.Network
	if lbNetwork != "" {
		if lbNetwork == "default" {
			_, err := client.GetDefaultNetwork()
			if err != nil {
				utility.Error("Error fetching default network: %s", err)
				os.Exit(1)
			}
		} else {
			_, err := client.FindNetwork(lbNetwork)
			if err != nil {
				utility.Error("Error finding network: %s", err)
				os.Exit(1)
			}
		}
		configLoadBalancer.NetworkID = lbNetwork
	}
	if lbCreateFirewall == "" {
		configLoadBalancer.FirewallRules = "80;443"
	} else {
		configLoadBalancer.FirewallRules = lbCreateFirewall
	}

	if lbExistingFirewall != "" {
		if lbCreateFirewall != "" {
			utility.Error("You can't use --create-firewall together with --existing-firewall flag")
			os.Exit(1)
		}

		ef, err := client.FindFirewall(lbExistingFirewall)
		if err != nil {
			utility.Error("Unable to find existing firewall %q: %s", lbExistingFirewall, err)
			os.Exit(1)
		}

		if ef.NetworkID != network.ID {
			utility.Error("Unable to find firewall %q in network %q", ef.ID, network.Label)
			os.Exit(1)
		}

		configLoadBalancer.FirewallID = ef.ID
		configLoadBalancer.FirewallRules = ""
	}
}

func setLoadBalancerOptions(configLoadBalancer *civogo.LoadBalancerConfig) {
	if lbAlgorithm != "" {
		configLoadBalancer.Algorithm = lbAlgorithm
	}

	if lbExternalTrafficPolicy != "" {
		configLoadBalancer.ExternalTrafficPolicy = lbExternalTrafficPolicy
	}

	if lbSessionAffinity != "" {
		configLoadBalancer.SessionAffinity = lbSessionAffinity
	}

	if lbSessionAffinityConfigTimeout != 0 {
		configLoadBalancer.SessionAffinityConfigTimeout = int32(lbSessionAffinityConfigTimeout)
	}

	if common.RegionSet != "" {
		configLoadBalancer.Region = common.RegionSet
	}
}

func setLoadBalancerBackends(configLoadBalancer *civogo.LoadBalancerConfig) error {
	var configLoadBalancerBackend []civogo.LoadBalancerBackendConfig

	for _, backend := range lbBackends {
		// Replace semicolons with colons to match GetStringMap's expected format
		backend = strings.ReplaceAll(backend, ";", ":")

		// Ensure that backend is non-empty
		if backend == "" {
			return fmt.Errorf("backend configuration cannot be empty")
		}

		data, _ := SetStringMap(backend)

		// Check if 'ip' is provided
		if ip, ok := data["ip"]; ok {
			backendConfig := civogo.LoadBalancerBackendConfig{
				IP: ip,
			}

			// Parse source-port
			if sourcePort, ok := data["source-port"]; ok {
				if port, err := strconv.Atoi(sourcePort); err == nil {
					backendConfig.SourcePort = int32(port)
				} else {
					return fmt.Errorf("invalid source-port: %s", err)
				}
			}

			// Parse target-port
			if targetPort, ok := data["target-port"]; ok {
				if port, err := strconv.Atoi(targetPort); err == nil {
					backendConfig.TargetPort = int32(port)
				} else {
					return fmt.Errorf("invalid target-port: %s", err)
				}
			}

			// Parse protocol
			if protocol, ok := data["protocol"]; ok {
				backendConfig.Protocol = protocol
			}

			// Parse health-check-port
			if healthCheckPort, ok := data["health-check-port"]; ok {
				if port, err := strconv.Atoi(healthCheckPort); err == nil {
					backendConfig.HealthCheckPort = int32(port)
				} else {
					return fmt.Errorf("invalid health-check-port: %s", err)
				}
			}

			configLoadBalancerBackend = append(configLoadBalancerBackend, backendConfig)
		} else {
			return fmt.Errorf("each backend must specify an 'ip' field.")
		}
	}

	configLoadBalancer.Backends = configLoadBalancerBackend
	return nil
}

func setLoadBalancerInstancePools(configLoadBalancer *civogo.LoadBalancerConfig) error {
	var configLoadBalancerInstancePools []civogo.LoadBalancerInstancePoolConfig

	for _, pool := range lbInstancePools {
		// Replace semicolons with colons for consistency with GetStringMap's expected format
		pool = strings.ReplaceAll(pool, ";", ":")

		// Ensure that pool is non-empty
		if pool == "" {
			return fmt.Errorf("instance pool configuration cannot be empty")
		}

		data, _ := SetStringMap(pool)

		instancePoolConfig := civogo.LoadBalancerInstancePoolConfig{}

		// Parse tags
		if tags, ok := data["tags"]; ok {
			instancePoolConfig.Tags = strings.Split(tags, ",")
		}

		// Parse names
		if names, ok := data["names"]; ok {
			instancePoolConfig.Names = strings.Split(names, ",")
		}

		// Parse protocol
		if protocol, ok := data["protocol"]; ok {
			instancePoolConfig.Protocol = protocol
		}

		// Parse source-port
		if sourcePort, ok := data["source-port"]; ok {
			if port, err := strconv.Atoi(sourcePort); err == nil {
				instancePoolConfig.SourcePort = int32(port)
			} else {
				return fmt.Errorf("invalid source-port: %s", err)
			}
		}

		// Parse target-port
		if targetPort, ok := data["target-port"]; ok {
			if port, err := strconv.Atoi(targetPort); err == nil {
				instancePoolConfig.TargetPort = int32(port)
			} else {
				return fmt.Errorf("invalid target-port: %s", err)
			}
		}

		// Parse health-check (port and path)
		instancePoolConfig.HealthCheck = civogo.HealthCheck{}
		if healthCheckPort, ok := data["health-check.port"]; ok {
			if port, err := strconv.Atoi(healthCheckPort); err == nil {
				instancePoolConfig.HealthCheck.Port = int32(port)
			} else {
				return fmt.Errorf("invalid health-check.port: %s", err)
			}
		}
		if healthCheckPath, ok := data["health-check.path"]; ok {
			instancePoolConfig.HealthCheck.Path = healthCheckPath
		}

		configLoadBalancerInstancePools = append(configLoadBalancerInstancePools, instancePoolConfig)
	}

	configLoadBalancer.InstancePools = configLoadBalancerInstancePools
	return nil
}

func outputLoadBalancer(loadBalancer *civogo.LoadBalancer) {
	ow := utility.NewOutputWriterWithMap(map[string]string{"id": loadBalancer.ID, "name": loadBalancer.Name})

	switch common.OutputFormat {
	case "json":
		ow.WriteSingleObjectJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		fmt.Printf("Created a new load balancer with name %s and ID %s\n", utility.Green(loadBalancer.Name), utility.Green(loadBalancer.ID))
	}
}
