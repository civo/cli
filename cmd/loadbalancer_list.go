package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var loadBalancerListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo loadbalancer ls -o custom -f "ID: Name"`,
	Short:   "List load balancers",
	Long: `List all load balancers.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Protocol
	* Port
	* TLSCertificate
	* TLSKey
	* Policy
	* HealthCheckPath
	* FailTimeout
	* MaxConns
	* IgnoreInvalidBackendTLS
	* Backends`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		lbs, err := client.ListLoadBalancers()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, lb := range lbs {
			ow.StartLine()

			ow.AppendData("ID", lb.ID)
			ow.AppendData("Name", lb.Hostname)
			ow.AppendData("Protocol", lb.Protocol)
			ow.AppendData("Port", strconv.Itoa(lb.Port))

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendDataWithLabel("TLSCertificate", lb.TLSCertificate, "TLS Cert")
				ow.AppendDataWithLabel("TLSKey", lb.TLSKey, "TLS Key")
				ow.AppendData("Policy", lb.Policy)
				ow.AppendDataWithLabel("HealthCheckPath", lb.HealthCheckPath, "Health Check Path")
				ow.AppendDataWithLabel("FailTimeout", strconv.Itoa(lb.FailTimeout), "Fail Timeout")
				ow.AppendDataWithLabel("MaxConns", strconv.Itoa(lb.MaxConns), "Max. Connections")
				ow.AppendDataWithLabel("IgnoreInvalidBackendTLS", strconv.FormatBool(lb.IgnoreInvalidBackendTLS), "Ignore Invalid Backend TLS?")
			}

			var backendList []string

			for _, backend := range lb.Backends {
				instance, err := client.FindInstance(backend.InstanceID)
				if err != nil {
					utility.Error("Finding the load balancer failed with %s", err)
					os.Exit(1)
				}
				backendList = append(backendList, instance.Hostname)
			}

			ow.AppendData("Backends", strings.Join(backendList, ", "))

		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
