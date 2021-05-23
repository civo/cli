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

	* id
	* name
	* protocol
	* port
	* tls_certificate
	* tls_key
	* policy
	* health_check_path
	* fail_timeout
	* max_conns
	* ignore_invalid_backend_tls`,
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

			ow.AppendDataWithLabel("id", lb.ID, "ID")
			ow.AppendDataWithLabel("name", lb.Hostname, "Name")
			ow.AppendDataWithLabel("protocol", lb.Protocol, "Protocol")
			ow.AppendDataWithLabel("port", strconv.Itoa(lb.Port), "Port")

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendDataWithLabel("tls_certificate", lb.TLSCertificate, "TLS Cert")
				ow.AppendDataWithLabel("tls_key", lb.TLSKey, "TLS Key")
				ow.AppendDataWithLabel("policy", lb.Policy, "Policy")
				ow.AppendDataWithLabel("health_check_path", lb.HealthCheckPath, "Health Check Path")
				ow.AppendDataWithLabel("fail_timeout", strconv.Itoa(lb.FailTimeout), "Fail Timeout")
				ow.AppendDataWithLabel("max_conns", strconv.Itoa(lb.MaxConns), "Max. Connections")
				ow.AppendDataWithLabel("ignore_invalid_backend_tls", strconv.FormatBool(lb.IgnoreInvalidBackendTLS), "Ignore Invalid Backend TLS?")
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
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
