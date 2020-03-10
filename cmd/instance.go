package cmd

import (
	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"instances"},
	Short:   "Details of Civo instances",
}

// TODO: instance create [--name=HOSTNAME] [...] -- create a new instance with specified hostname and provided options
// TODO: instance tags ID/HOSTNAME 'tag1 tag2 tag3...' -- retag instance by ID (input no tags to clear all tags)
// TODO: instance update ID/HOSTNAME [--name] [--notes] -- update details of instance

func init() {
	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceShowCmd)
	instanceCmd.AddCommand(instanceRemoveCmd)
	instanceCmd.AddCommand(instanceRebootCmd)
	instanceCmd.AddCommand(instanceSoftRebootCmd)
	instanceCmd.AddCommand(instanceConsoleCmd)
	instanceCmd.AddCommand(instanceStopCmd)
	instanceCmd.AddCommand(instanceStartCmd)
	instanceCmd.AddCommand(instanceUpgradeCmd)
	instanceCmd.AddCommand(instanceMoveIPCmd)
	instanceCmd.AddCommand(instanceSetFirewallCmd)
	instanceCmd.AddCommand(instancePublicIPCmd)
	instanceCmd.AddCommand(instancePasswordCmd)
}
