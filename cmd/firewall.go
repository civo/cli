package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
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
	rootCmd.AddCommand(firewallCmd)
	firewallCmd.AddCommand(firewallListCmd)
	firewallCmd.AddCommand(firewallCreateCmd)
	firewallCmd.AddCommand(firewallUpdateCmd)
	firewallCmd.AddCommand(firewallRemoveCmd)

	firewallCreateCmd.Flags().StringVarP(&firewallnetwork, "network", "n", "default", "the network to create the firewall")

	// Firewalls rule cmd
	firewallCmd.AddCommand(firewallRuleCmd)
	firewallRuleCmd.AddCommand(firewallRuleListCmd)
	firewallRuleCmd.AddCommand(firewallRuleCreateCmd)
	firewallRuleCmd.AddCommand(firewallRuleRemoveCmd)

	/*
		Flags for firewall rule create cmd
	*/
	firewallRuleCreateCmd.Flags().StringVarP(&protocol, "protocol", "p", "", "the protocol choice (from: TCP, UDP, ICMP)")
	firewallRuleCreateCmd.Flags().StringVarP(&startPort, "startport", "s", "", "the start port of the rule")
	firewallRuleCreateCmd.Flags().StringVarP(&endPort, "endport", "e", "", "the end port of the rule")
	firewallRuleCreateCmd.Flags().StringArrayVarP(&cidr, "cidr", "c", []string{}, "the CIDR of the rule you can use (e.g. -c 10.10.10.1/32, 10.10.10.2/32)")
	firewallRuleCreateCmd.Flags().StringVarP(&direction, "direction", "d", "", "the direction of the rule need to be ingress")
	firewallRuleCreateCmd.Flags().StringVarP(&label, "label", "l", "", "a string that will be the displayed as the name/reference for this rule")
}
