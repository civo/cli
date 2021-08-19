package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"instances"},
	Short:   "Details of Civo instances",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceSizeCmd)
	instanceCmd.AddCommand(instanceCreateCmd)
	instanceCmd.AddCommand(instanceShowCmd)
	instanceCmd.AddCommand(instanceUpdateCmd)
	instanceCmd.AddCommand(instanceRemoveCmd)
	instanceCmd.AddCommand(instanceRebootCmd)
	instanceCmd.AddCommand(instanceSoftRebootCmd)
	instanceCmd.AddCommand(instanceConsoleCmd)
	instanceCmd.AddCommand(instanceStopCmd)
	instanceCmd.AddCommand(instanceStartCmd)
	instanceCmd.AddCommand(instanceUpgradeCmd)
	// instanceCmd.AddCommand(instanceMoveIPCmd)
	instanceCmd.AddCommand(instanceSetFirewallCmd)
	instanceCmd.AddCommand(instancePublicIPCmd)
	instanceCmd.AddCommand(instancePasswordCmd)
	instanceCmd.AddCommand(instanceTagCmd)

	instanceUpdateCmd.Flags().StringVarP(&notes, "notes", "n", "", "notes stored against the instance")
	instanceUpdateCmd.Flags().StringVarP(&reverseDNS, "reverse-dns", "r", "", "the reverse DNS entry for the instance")
	instanceUpdateCmd.Flags().StringVarP(&hostname, "hostname", "s", "", "the instance's hostname")

	instanceCreateCmd.Flags().BoolVarP(&wait, "wait", "w", false, "wait until the instance's is ready")
	instanceCreateCmd.Flags().StringVarP(&hostnameCreate, "hostname", "s", "", "the instance's hostname")
	instanceCreateCmd.Flags().StringVarP(&size, "size", "i", "", "the instance's size (from 'civo instance size' command)")
	instanceCreateCmd.MarkFlagRequired("size")
	instanceCreateCmd.Flags().StringVarP(&template, "template", "t", "", "the instance's template (from 'civo template ls' command)")
	instanceCreateCmd.MarkFlagRequired("template")
	instanceCreateCmd.Flags().StringVarP(&snapshot, "snapshot", "n", "", "the instance's snapshot")
	instanceCreateCmd.Flags().StringVarP(&publicip, "publicip", "p", "create", "This should be either none, create or `move_ip_from:intances_id` by default is create")
	instanceCreateCmd.Flags().StringVarP(&initialuser, "initialuser", "u", "", "the instance's initial user")
	instanceCreateCmd.Flags().StringVarP(&sshkey, "sshkey", "k", "", "the instance's ssh key you can use the Name or the ID")
	instanceCreateCmd.Flags().StringVarP(&network, "network", "r", "", "the instance's network you can use the Name or the ID")
	instanceCreateCmd.Flags().StringVarP(&tags, "tags", "g", "", "the instance's tags")
	instanceCreateCmd.Flags().StringVarP(&tags, "region", "e", "", "the region code identifier to have your instance built in")

	instanceStopCmd.Flags().BoolVarP(&waitStop, "wait", "w", false, "wait until the instance's is stoped")
}
