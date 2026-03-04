package vpc

import (
	"errors"

	"github.com/spf13/cobra"
)

var vpcLoadBalancerCmd = &cobra.Command{
	Use:     "loadbalancer",
	Aliases: []string{"loadbalancers", "lb"},
	Short:   "Details of Civo VPC load balancers",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}
