package vpc

import (
	"errors"

	"github.com/spf13/cobra"
)

var vpcFirewallRuleCmd = &cobra.Command{
	Use:     "rule",
	Aliases: []string{"rules"},
	Short:   "Details of Civo VPC firewall rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}
