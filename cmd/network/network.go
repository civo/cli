package network

import (
	"errors"

	"github.com/spf13/cobra"
)

var v4enabled, v6enabled bool
var cidrv4 string

// NetworkCmd manages Civo networks
var NetworkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"networks", "net"},
	Short:   "Details of Civo networks",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var networkSubnetCmd = &cobra.Command{
	Use:     "subnet",
	Aliases: []string{"subnets"},
	Short:   "Details of Civo network subnets",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	NetworkCmd.AddCommand(networkListCmd)
	NetworkCmd.AddCommand(networkCreateCmd)
	NetworkCmd.AddCommand(networkUpdateCmd)
	NetworkCmd.AddCommand(networkRemoveCmd)
	NetworkCmd.AddCommand(networkSubnetCmd)

	networkCreateCmd.Flags().BoolVarP(&v4enabled, "v4", "", true, "Enable IPv4 on the network")
	networkCreateCmd.Flags().BoolVarP(&v6enabled, "v6", "", false, "Enable IPv6 on the network")
	networkCreateCmd.Flags().StringVarP(&cidrv4, "cidr-v4", "c", "", "The CIDR for the IPv4 network")

	networkSubnetCmd.AddCommand(networkSubnetListCmd)
	networkSubnetCmd.AddCommand(networkSubnetCreateCmd)
	networkSubnetCmd.AddCommand(networkSubnetShowCmd)
	networkSubnetCmd.AddCommand(networkSubnetRemoveCmd)
	networkSubnetCmd.AddCommand(networkSubnetAttachCmd)
	networkSubnetCmd.AddCommand(networkSubnetDetachCmd)

	networkSubnetAttachCmd.Flags().StringVarP(&subnetID, "subnet", "", "", "the id of subnet you want to attach to the instance")
	networkSubnetAttachCmd.Flags().StringVarP(&instanceID, "instance", "", "", "the id of instance you want to attach your subnet to")
	networkSubnetAttachCmd.MarkFlagRequired("subnet")
	networkSubnetAttachCmd.MarkFlagRequired("instance")

	networkSubnetDetachCmd.Flags().StringVarP(&subnetID, "subnet", "", "", "the id of subnet you want to attach to the instance")
	networkSubnetDetachCmd.Flags().StringVarP(&instanceID, "instance", "", "", "the id of instance you want to attach your subnet to")
	networkSubnetDetachCmd.MarkFlagRequired("subnet")
	networkSubnetDetachCmd.MarkFlagRequired("instance")
}
