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

func init() {
	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceShowCmd)
	instanceCmd.AddCommand(instanceUpdateCmd)
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

	instanceCmd.AddCommand(instanceTagCmd)
	instanceUpdateCmd.Flags().StringVarP(&notes, "notes", "n", "", "notes stored against the instance")
	instanceUpdateCmd.Flags().StringVarP(&reverseDNS, "reverse-dns", "r", "", "the reverse DNS entry for the instance")
	instanceUpdateCmd.Flags().StringVarP(&hostname, "hostname", "e", "", "the instance's hostname")
}
