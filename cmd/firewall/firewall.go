package firewall

import (
	"errors"

	"github.com/spf13/cobra"
)

// FirewallCmd manages Civo firewalls
var FirewallCmd = &cobra.Command{
	Use:     "firewall",
	Aliases: []string{"firewalls", "fw"},
	Short:   "Details of Civo firewalls",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var firewallRuleCmd = &cobra.Command{
	Use:     "rule",
	Aliases: []string{"rules"},
	Short:   "Details of Civo firewalls rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	FirewallCmd.AddCommand(firewallListCmd)
	FirewallCmd.AddCommand(firewallCreateCmd)
	FirewallCmd.AddCommand(firewallUpdateCmd)
	FirewallCmd.AddCommand(firewallRemoveCmd)

	firewallCreateCmd.Flags().StringVarP(&firewallnetwork, "network", "n", "default", "the network to create the firewall")
	firewallCreateCmd.Flags().BoolVarP(&createRules, "create-rules", "r", true, "the create rules flag is used to create the default firewall rules, if is not defined will be set to true (deprecated)")
	firewallCreateCmd.Flags().BoolVarP(&noDefaultRules, "no-default-rules", "", false, "the no-default-rules flag will ensure no default rules are created for the firewall, if not defined it will be set to false")

	// Firewalls rule cmd
	FirewallCmd.AddCommand(firewallRuleCmd)
	firewallRuleCmd.AddCommand(firewallRuleListCmd)
	firewallRuleCmd.AddCommand(firewallRuleCreateCmd)
	firewallRuleCmd.AddCommand(firewallRuleRemoveCmd)

	// Flags for firewall rule create cmd
	firewallRuleCreateCmd.Flags().StringVarP(&protocol, "protocol", "p", "TCP", "the protocol choice (TCP, UDP, ICMP)")
	firewallRuleCreateCmd.Flags().StringVarP(&startPort, "startport", "s", "", "the start port of the rule")
	firewallRuleCreateCmd.Flags().StringVarP(&endPort, "endport", "e", "", "the end port of the rule")
	firewallRuleCreateCmd.Flags().StringVarP(&cidr, "cidr", "c", "0.0.0.0/0", "the CIDR of the rule you can use (e.g. -c 10.10.10.1/32,148.2.6.120/32)")
	firewallRuleCreateCmd.Flags().StringVarP(&direction, "direction", "d", "ingress", "the direction of the rule can be ingress or egress (default is ingress)")
	firewallRuleCreateCmd.Flags().StringVarP(&action, "action", "a", "allow", "the action of the rule can be allow or deny (default is allow)")
	firewallRuleCreateCmd.Flags().StringVarP(&label, "label", "l", "", "a string that will be the displayed as the name/reference for this rule")
	_ = firewallRuleCreateCmd.MarkFlagRequired("startport")

	// Mark the create-rules flag as deprecated
	_ = firewallCreateCmd.Flags().MarkDeprecated("create-rules", "it will be removed in future versions. Default firewall rules are created by default. Use --no-default-rules flag to create firewalls without them.\n")
}
