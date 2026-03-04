package vpc

import (
	"errors"

	"github.com/spf13/cobra"
)

var vpcSubnetCmd = &cobra.Command{
	Use:     "subnet",
	Aliases: []string{"subnets"},
	Short:   "Details of Civo VPC subnets",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}
