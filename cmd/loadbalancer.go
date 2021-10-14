package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var loadBalancerCmd = &cobra.Command{
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
	// rootCmd.AddCommand(loadBalancerCmd)
	loadBalancerCmd.AddCommand(loadBalancerListCmd)
	loadBalancerCmd.AddCommand(loadBalancerRemoveCmd)
	loadBalancerCmd.AddCommand(loadBalancerCreateCmd)
	loadBalancerCmd.AddCommand(loadBalancerUpdateCmd)

	loadBalancerCreateCmd.Flags().StringVarP(&lbHostname, "hostname", "e", "", "If not supplied, will be in format loadbalancer-uuid.civo.com")
	loadBalancerCreateCmd.Flags().StringVarP(&lbProtocol, "protocol", "p", "", "Either http or https. If you specify https then you must also provide the next two fields")
	loadBalancerCreateCmd.Flags().StringVarP(&tlsCertificate, "tls_certificate", "c", "", "TLS certificate in Base64-encoded PEM. Required if --protocol is https")
	loadBalancerCreateCmd.Flags().StringVarP(&tlsKey, "tls_key", "k", "", "TLS certificate in Base64-encoded PEM. Required if --protocol is https")
	loadBalancerCreateCmd.Flags().StringVarP(&policy, "policy", "", "", "<least_conn | random | round_robin | ip_hash> - Balancing policy to choose backends")
	loadBalancerCreateCmd.Flags().IntVarP(&lbPort, "port", "r", 80, "Listening port. Defaults to 80 to match default http protocol")
	loadBalancerCreateCmd.Flags().IntVarP(&maxRequestSize, "max_request_size", "m", 20, "Maximum request content size, in MB. Defaults to 20")
	loadBalancerCreateCmd.Flags().StringVarP(&healthCheckPath, "health_check_path", "l", "", "URL to check for a valid (2xx/3xx) HTTP status on the backends. Defaults to /")
	loadBalancerCreateCmd.Flags().IntVarP(&failTimeout, "fail_timeout", "t", 30, "Timeout in seconds to consider a backend to have failed. Defaults to 30")
	loadBalancerCreateCmd.Flags().IntVarP(&maxConnections, "max_connections", "x", 10, "Maximum concurrent connections to each backend. Defaults to 10")
	loadBalancerCreateCmd.Flags().BoolVarP(&ignoreInvalidBackendTLS, "ignore_invalid_backend_tls", "i", true, "Should self-signed/invalid certificates be ignored from backend servers? Defaults to true")
	loadBalancerCreateCmd.Flags().StringArrayVarP(&backends, "backends", "b", []string{}, "Specify a backend instance to associate with the load balancer. Takes instance_id, protocol and port in the format --backend=instance:instance-id|instance-name,protocol:http,port:80")

	loadBalancerUpdateCmd.Flags().StringVarP(&lbHostnameUpdate, "hostname", "e", "", "If not supplied, will be in format loadbalancer-uuid.civo.com")
	loadBalancerUpdateCmd.Flags().StringVarP(&lbProtocolUpdate, "protocol", "p", "", "Either http or https. If you specify https then you must also provide the next two fields")
	loadBalancerUpdateCmd.Flags().StringVarP(&tlsCertificateUpdate, "tls_certificate", "c", "", "TLS certificate in Base64-encoded PEM. Required if --protocol is https")
	loadBalancerUpdateCmd.Flags().StringVarP(&tlsKeyUpdate, "tls_key", "k", "", "TLS certificate in Base64-encoded PEM. Required if --protocol is https")
	loadBalancerUpdateCmd.Flags().StringVarP(&policyUpdate, "policy", "", "", "<least_conn | random | round_robin | ip_hash> - Balancing policy to choose backends")
	loadBalancerUpdateCmd.Flags().IntVarP(&lbPortUpdate, "port", "r", 80, "Listening port. Defaults to 80 to match default http protocol")
	loadBalancerUpdateCmd.Flags().IntVarP(&maxRequestSizeUpdate, "max_request_size", "m", 20, "Maximum request content size, in MB. Defaults to 20")
	loadBalancerUpdateCmd.Flags().StringVarP(&healthCheckPathUpdate, "health_check_path", "l", "", "URL to check for a valid (2xx/3xx) HTTP status on the backends. Defaults to /")
	loadBalancerUpdateCmd.Flags().IntVarP(&failTimeoutUpdate, "fail_timeout", "t", 30, "Timeout in seconds to consider a backend to have failed. Defaults to 30")
	loadBalancerUpdateCmd.Flags().IntVarP(&maxConnectionsUpdate, "max_connections", "x", 10, "Maximum concurrent connections to each backend. Defaults to 10")
	loadBalancerUpdateCmd.Flags().BoolVarP(&ignoreInvalidBackendTLSUpdate, "ignore_invalid_backend_tls", "i", true, "Should self-signed/invalid certificates be ignored from backend servers? Defaults to true")
	loadBalancerUpdateCmd.Flags().StringArrayVarP(&backendsUpdate, "backends", "b", []string{}, "Specify a backend instance to associate with the load balancer. Takes instance_id, protocol and port in the format --backend=instance:instance-id|instance-name,protocol:http,port:80")

}
