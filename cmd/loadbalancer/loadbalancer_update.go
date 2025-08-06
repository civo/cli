package loadbalancer

import (
	"fmt"
	"github.com/civo/cli/common"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var lbNameUpdate, lbAlgorithmUpdate string
var lbBackendsUpdate []string
var lbInstancePoolsUpdate []string

// New variables for instance pool updates
// var instancePoolTags []string
// var instancePoolNames []string
// var instancePoolProtocol string
var instancePoolSourcePort, instancePoolTargetPort int32
var instancePoolHealthCheckPort int32
var instancePoolHealthCheckPath string

var loadBalancerUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Example: "civo loadbalancer update ID/NAME [flags]",
	Short:   "Update a load balancer",
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

		loadBalancer, err := client.FindLoadBalancer(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		configLoadBalancer := &civogo.LoadBalancerUpdateConfig{Region: client.Region}

		if lbNameUpdate != "" {
			configLoadBalancer.Name = lbNameUpdate
		}

		if lbAlgorithmUpdate != "" {
			configLoadBalancer.Algorithm = lbAlgorithmUpdate
		}

		// Handle backend updates
		if len(lbBackendsUpdate) > 0 {
			var configLoadBalancerBackends []civogo.LoadBalancerBackendConfig
			for _, backend := range lbBackendsUpdate {
				// setStringMap converts a semicolon-separated string into a map, supporting nested keys like "health-check.port"
				data, err := SetStringMap(backend)
				if err != nil {
					utility.Error("Error parsing backend data: %s", err)
				}

				loadBalancerBackend := civogo.LoadBalancerBackendConfig{
					IP: data["ip"],
				}

				if data["protocol"] != "" {
					loadBalancerBackend.Protocol = data["protocol"]
				}

				if data["source-port"] != "" {
					if port, err := strconv.Atoi(data["source-port"]); err == nil {
						loadBalancerBackend.SourcePort = int32(port)
					}
				}

				if data["target-port"] != "" {
					if port, err := strconv.Atoi(data["target-port"]); err == nil {
						loadBalancerBackend.TargetPort = int32(port)
					}
				}

				if data["health-check-port"] != "" {
					if port, err := strconv.Atoi(data["health-check-port"]); err == nil {
						loadBalancerBackend.HealthCheckPort = int32(port)
					}
				}

				configLoadBalancerBackends = append(configLoadBalancerBackends, loadBalancerBackend)
			}
			configLoadBalancer.Backends = configLoadBalancerBackends
		}

		// Handle instance pool updates
		if len(lbInstancePoolsUpdate) > 0 {
			var instancePools []civogo.LoadBalancerInstancePoolConfig
			for _, pool := range lbInstancePoolsUpdate {
				data, err := SetStringMap(pool)
				if err != nil {
					utility.Error("Error parsing instance pool data: %s", err)
				}

				instancePool := civogo.LoadBalancerInstancePoolConfig{
					SourcePort: instancePoolSourcePort,
					TargetPort: instancePoolTargetPort,
					HealthCheck: civogo.HealthCheck{
						Port: instancePoolHealthCheckPort,
						Path: instancePoolHealthCheckPath,
					},
				}

				if data["tags"] != "" {
					instancePool.Tags = utility.SplitCommaSeparatedValues(data["tags"])
				}

				if data["names"] != "" {
					instancePool.Names = utility.SplitCommaSeparatedValues(data["names"])
				}

				// Update instance pool only if values are provided
				if data["protocol"] != "" {
					instancePool.Protocol = data["protocol"]
				}
				if data["source-port"] != "" {
					if port, err := strconv.Atoi(data["source-port"]); err == nil {
						instancePool.SourcePort = int32(port)
					}
				}
				if data["target-port"] != "" {
					if port, err := strconv.Atoi(data["target-port"]); err == nil {
						instancePool.TargetPort = int32(port)
					}
				}
				if data["health-check-port"] != "" {
					if port, err := strconv.Atoi(data["health-check-port"]); err == nil {
						instancePool.HealthCheck.Port = int32(port)
					}
				}
				if data["health-check-path"] != "" {
					instancePool.HealthCheck.Path = data["health-check-path"]
				}

				instancePools = append(instancePools, instancePool)
			}
			configLoadBalancer.InstancePools = instancePools
		}

		loadBalancerUpdate, err := client.UpdateLoadBalancer(loadBalancer.ID, configLoadBalancer)
		if err != nil {
			utility.Error("error while updating the LB: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": loadBalancerUpdate.ID, "hostname": loadBalancerUpdate.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Updated load balancer with name %s with ID %s\n", utility.Green(loadBalancerUpdate.Name), utility.Green(loadBalancerUpdate.ID))
		}
	},
}

// SetStringMap converts a semicolon-separated string into a map, supporting nested keys like "health-check.port"
func SetStringMap(input string) (map[string]string, error) {
	result := make(map[string]string)
	if input == "" {
		return result, nil
	}

	// Split by semicolon
	entries := strings.Split(input, ";")
	for _, entry := range entries {
		// Split by equal sign
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			// Log the error and skip the malformed entry
			fmt.Printf("Warning: Skipping invalid entry format: '%s'\n", entry)
			continue
		}
		// Add key-value pair to the result map
		result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return result, nil
}
