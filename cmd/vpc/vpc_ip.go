package vpc

import (
	"errors"

	"github.com/spf13/cobra"
)

var vpcIPCmd = &cobra.Command{
	Use:     "ip",
	Aliases: []string{"ips"},
	Short:   "Details of Civo VPC reserved IPs",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}
