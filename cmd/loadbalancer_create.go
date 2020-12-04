package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var lbHostname, lbProtocol, tlsCertificate, tlsKey, policy, healthCheckPath string
var lbPort, maxRequestSize, failTimeout, maxConnections int
var ignoreInvalidBackendTLS bool
var backends []string

var loadBalancerCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo loadbalancer create [flags]",
	Short:   "Create a new load balancer",
	//Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
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
					utility.Error("%s", err)
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
			utility.Error("Creating the load balancer failed with %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": loadBalancer.ID, "Hostname": loadBalancer.Hostname})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a new load balancer with hostname %s with ID %s\n", utility.Green(loadBalancer.Hostname), utility.Green(loadBalancer.ID))
		}
	},
}
