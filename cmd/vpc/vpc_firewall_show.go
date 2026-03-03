package vpc

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcFirewallShowCmd = &cobra.Command{
	Use:     "show [FIREWALL-NAME/FIREWALL-ID]",
	Short:   "Show details of a specific VPC firewall",
	Aliases: []string{"get", "describe", "inspect"},
	Args:    cobra.ExactArgs(1),
	Example: "civo vpc firewall show FIREWALL_NAME",
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

		firewall, err := client.FindVPCFirewall(args[0])
		if err != nil {
			utility.Error("Firewall %s", err)
			os.Exit(1)
		}

		rules, err := client.ListVPCFirewallRules(firewall.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", firewall.ID, "ID")
		ow.AppendDataWithLabel("name", firewall.Name, "Name")
		ow.AppendDataWithLabel("network_id", firewall.NetworkID, "Network ID")
		ow.AppendDataWithLabel("rules_count", fmt.Sprintf("%d", firewall.RulesCount), "Rules Count")

		switch common.OutputFormat {
		case "json":
			ow.ToJSON(firewall, common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Println("VPC Firewall Details:")
			fmt.Printf("ID: %s\n", firewall.ID)
			fmt.Printf("Name: %s\n", firewall.Name)
			fmt.Printf("Network ID: %s\n", firewall.NetworkID)
			fmt.Printf("Rules Count: %d\n", firewall.RulesCount)
			fmt.Printf("Instance Count: %d\n", firewall.InstanceCount)
			fmt.Printf("Cluster Count: %d\n", firewall.ClusterCount)
			fmt.Printf("Load Balancer Count: %d\n", firewall.LoadBalancerCount)

			if len(rules) > 0 {
				fmt.Println("\nRules:")
				for _, rule := range rules {
					fmt.Printf("  - %s: %s %s ports %s-%s %s (%s)\n",
						rule.ID, rule.Direction, rule.Protocol,
						rule.StartPort, rule.EndPort,
						strings.Join(rule.Cidr, ","), rule.Action)
				}
			}
		}
	},
}
