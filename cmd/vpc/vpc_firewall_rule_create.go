package vpc

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	vpcFwRuleProtocol  string
	vpcFwRuleStartPort string
	vpcFwRuleEndPort   string
	vpcFwRuleCidr      string
	vpcFwRuleDirection string
	vpcFwRuleAction    string
	vpcFwRuleLabel     string
)

var vpcFirewallRuleCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new VPC firewall rule",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo vpc firewall rule create FIREWALL_NAME/FIREWALL_ID [flags]",
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
			utility.Error("%s", err)
			os.Exit(1)
		}

		if err := vpcValidateCIDRs(vpcFwRuleCidr); err != nil {
			utility.Error("%s", err.Error())
			os.Exit(1)
		}

		newRuleConfig := &civogo.FirewallRuleConfig{
			FirewallID: firewall.ID,
			Protocol:   vpcFwRuleProtocol,
			StartPort:  vpcFwRuleStartPort,
			Cidr:       strings.Split(vpcFwRuleCidr, ","),
			Label:      vpcFwRuleLabel,
			Action:     vpcFwRuleAction,
			Region:     client.Region,
		}

		var directionValue string
		switch vpcFwRuleDirection {
		case "ingress":
			newRuleConfig.Direction = vpcFwRuleDirection
			directionValue = "from"
		case "egress":
			newRuleConfig.Direction = vpcFwRuleDirection
			directionValue = "to"
		default:
			utility.Error("'--direction' flag must be 'ingress' or 'egress'")
			os.Exit(1)
		}

		if vpcFwRuleEndPort == "" {
			newRuleConfig.EndPort = vpcFwRuleStartPort
		} else {
			newRuleConfig.EndPort = vpcFwRuleEndPort
		}

		rule, err := client.NewVPCFirewallRule(newRuleConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": rule.ID, "name": rule.Label})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if rule.Label == "" {
				fmt.Printf("Created a VPC firewall rule to %s, %s access to port %s %s %s with ID %s\n",
					utility.Green(newRuleConfig.Action), utility.Green(newRuleConfig.Direction),
					utility.Green(newRuleConfig.StartPort), directionValue,
					utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
			} else {
				fmt.Printf("Created a VPC firewall rule called %s to %s, %s access to port %s %s %s with ID %s\n",
					utility.Green(rule.Label), utility.Green(newRuleConfig.Action),
					utility.Green(newRuleConfig.Direction), utility.Green(newRuleConfig.StartPort),
					directionValue, utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
			}
		}
	},
}

func vpcValidateCIDRs(cidrs string) error {
	for _, cidr := range strings.Split(cidrs, ",") {
		if cidr = strings.TrimSpace(cidr); cidr == "" {
			continue
		}
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return fmt.Errorf("invalid CIDR address '%s': %s", cidr, err)
		}
	}
	return nil
}
