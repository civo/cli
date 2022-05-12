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
	// instanceCmd.AddCommand(instanceConsoleCmd)
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
	instanceCreateCmd.Flags().StringVarP(&diskimage, "diskimage", "t", "", "the instance's disk image (from 'civo diskimage ls' command)")
	instanceCreateCmd.Flags().StringVarP(&publicip, "publicip", "p", "create", "This should be either \"none\" or \"create\"")
	instanceCreateCmd.Flags().StringVarP(&initialuser, "initialuser", "u", "", "the instance's initial user")
	instanceCreateCmd.Flags().StringVarP(&sshkey, "sshkey", "k", "", "the instance's ssh key you can use the Name or the ID")
	instanceCreateCmd.Flags().StringVarP(&network, "network", "r", "", "the instance's network you can use the Name or the ID")
	instanceCreateCmd.Flags().StringVarP(&firewall, "firewall", "l", "", "the instance's firewall you can use the Name or the ID")
	instanceCreateCmd.Flags().StringVarP(&tags, "tags", "g", "", "the instance's tags")
	instanceCreateCmd.Flags().StringVarP(&tags, "region", "e", "", "the region code identifier to have your instance built in")
	instanceCreateCmd.Flags().StringVar(&script, "script", "", "path to a script that will be uploaded to /usr/local/bin/civo-user-init-script on your instance, read/write/executable only by root and then will be executed at the end of the cloud initialization")
	instanceCreateCmd.Flags().BoolVar(&skipShebangCheck, "skip-shebang-check", false, "skip the shebang line check when passing a user init script")

	instanceStopCmd.Flags().BoolVarP(&waitStop, "wait", "w", false, "wait until the instance's is stoped")
}
