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

var lbHostnameUpdate, lbProtocolUpdate, tlsCertificateUpdate, tlsKeyUpdate, policyUpdate, healthCheckPathUpdate string
var lbPortUpdate, maxRequestSizeUpdate, failTimeoutUpdate, maxConnectionsUpdate int
var ignoreInvalidBackendTLSUpdate bool
var backendsUpdate []string

var loadBalancerUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Example: "civo loadbalancer update ID/HOSTNAME [flags]",
	Short:   "Update a load balancer",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		loadBalancer, err := client.FindLoadBalancer(args[0])
		if err != nil {
			utility.Error("%s", err)
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
					utility.Error("Instance %s", err)
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
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": loadBalancerUpdate.ID, "hostname": loadBalancerUpdate.Hostname})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Updated load balancer with hostname %s with ID %s\n", utility.Green(loadBalancerUpdate.Hostname), utility.Green(loadBalancerUpdate.ID))
		}
	},
}
