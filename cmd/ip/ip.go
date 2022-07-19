package ip

import (
	"errors"

	"github.com/spf13/cobra"
)

var name, instance string

// IPCmd manages reserved IPs
var IPCmd = &cobra.Command{
	Use:     "ip",
	Aliases: []string{"ips"},
	Short:   "Details of Civo reserved IPs",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	IPCmd.AddCommand(ipListCmd)
	IPCmd.AddCommand(ipCreateCmd)
	IPCmd.AddCommand(ipRenameCmd)
	IPCmd.AddCommand(ipDeleteCmd)
	IPCmd.AddCommand(ipAssignCmd)
	IPCmd.AddCommand(ipUnassignCmd)

	ipCreateCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the reserved IP")
	ipAssignCmd.Flags().StringVarP(&instance, "instance", "i", "", "Name/ID of the instance to assign the IP to")
	ipAssignCmd.MarkFlagRequired("instance")
}
