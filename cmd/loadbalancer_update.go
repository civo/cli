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

var lbHostnameUpdate, lbProtocolUpdate, tlsCertificateUpdate, tlsKeyUpdate, policyUpdate, healthCheckPathUpdate string
var lbPortUpdate, maxRequestSizeUpdate, failTimeoutUpdate, maxConnectionsUpdate int
var ignoreInvalidBackendTLSUpdate bool
var backendsUpdate []string

var loadBalancerUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Example: "civo loadbalancer update [flags]",
	Short:   "Update a Load Balancer",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		loadBalancer, err := client.FindLoadBalancer(args[0])
		if err != nil {
			utility.Error("Unable to find the load balancer %s", err)
			os.Exit(1)
		}

		configLoadBalancer := &civogo.LoadBalancerConfig{
			Port:                    lbPortUpdate,
			MaxRequestSize:          maxRequestSizeUpdate,
			FailTimeout:             failTimeoutUpdate,
			MaxConns:                maxConnectionsUpdate,
			IgnoreInvalidBackendTLS: ignoreInvalidBackendTLSUpdate,
		}

		if lbHostnameUpdate != "" {
			configLoadBalancer.Hostname = lbHostnameUpdate
		} else {
			configLoadBalancer.Hostname = loadBalancer.Hostname
		}

		if lbProtocolUpdate != "" {
			configLoadBalancer.Protocol = lbProtocolUpdate
		}

		if tlsCertificateUpdate != "" {
			configLoadBalancer.TLSCertificate = tlsCertificateUpdate
		}

		if tlsKeyUpdate != "" {
			configLoadBalancer.TLSKey = tlsKeyUpdate
		}

		if policyUpdate != "" {
			configLoadBalancer.Policy = policyUpdate
		}

		if healthCheckPathUpdate != "" {
			configLoadBalancer.HealthCheckPath = healthCheckPathUpdate
		}

		if len(backendsUpdate) > 0 {
			var configLoadBalancerBackend []civogo.LoadBalancerBackendConfig

			for _, backend := range backendsUpdate {
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

		loadBalancerUpdate, err := client.UpdateLoadBalancer(loadBalancer.ID, configLoadBalancer)
		if err != nil {
			utility.Error("Unable to update a load balancer %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": loadBalancerUpdate.ID, "Hostname": loadBalancerUpdate.Hostname})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Updated Load Balancer with hostname %s with ID %s\n", utility.Green(loadBalancerUpdate.Hostname), utility.Green(loadBalancerUpdate.ID))
		}
	},
}
