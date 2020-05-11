package cmd

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var lbHostname, lbProtocol, tlsCertificate, tlsKey, policy, healthCheckPath string
var lbPort, maxRequestSize, failTimeout, maxConnections int
var ignoreInvalidBackendTLS bool
var backends []string

var loadBalancerCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new Load Balancer",
	//Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		configLoadBalancer := &civogo.LoadBalancerConfig{
			Port:                    lbPort,
			MaxRequestSize:          maxRequestSize,
			FailTimeout:             failTimeout,
			MaxConns:                maxConnections,
			IgnoreInvalidBackendTLS: ignoreInvalidBackendTLS,
		}

		if lbHostname != "" {
			configLoadBalancer.Hostname = lbHostname
		}

		if lbProtocol != "" {
			configLoadBalancer.Protocol = lbProtocol
		}

		if tlsCertificate != "" {
			configLoadBalancer.TLSCertificate = tlsCertificate
		}

		if tlsKey != "" {
			configLoadBalancer.TLSKey = tlsKey
		}

		if policy != "" {
			configLoadBalancer.Policy = policy
		}

		if healthCheckPath != "" {
			configLoadBalancer.HealthCheckPath = healthCheckPath
		}

		if len(backends) > 0 {
			var configLoadBalancerBackend []civogo.LoadBalancerBackendConfig

			for _, backend := range backends {
				data := utility.GetStringMap(backend)
				instance, err := client.FindInstance(data["instance"])
				if err != nil {
					utility.Error("Unable to find the instance %s", err)
					os.Exit(1)
				}

				portBackend, err := strconv.Atoi(data["port"])
				if err != nil {
					fmt.Println(err)
				}

				configLoadBalancerBackend = append(configLoadBalancerBackend, civogo.LoadBalancerBackendConfig{
					InstanceID: instance.ID,
					Protocol:   data["protocol"],
					Port:       portBackend,
				})
			}

			configLoadBalancer.Backends = configLoadBalancerBackend
		}

		loadBalancer, err := client.CreateLoadBalancer(configLoadBalancer)
		if err != nil {
			utility.Error("Unable to create a load balancer %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": loadBalancer.ID, "Hostname": loadBalancer.Hostname})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a new Load Balancer with hostname %s with ID %s\n", utility.Green(loadBalancer.Hostname), utility.Green(loadBalancer.ID))
		}
	},
}
