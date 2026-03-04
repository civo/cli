package vpc

import (
	"errors"

	"github.com/spf13/cobra"
)

var vpcFirewallCmd = &cobra.Command{
	Use:     "firewall",
	Aliases: []string{"firewalls", "fw"},
	Short:   "Details of Civo VPC firewalls",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}
