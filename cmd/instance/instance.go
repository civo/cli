package instance

import (
	"errors"

	"github.com/spf13/cobra"
)

// InstanceCmd manages Civo instances
var InstanceCmd = &cobra.Command{
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
	InstanceCmd.AddCommand(instanceListCmd)
	InstanceCmd.AddCommand(instanceSizeCmd)
	InstanceCmd.AddCommand(instanceCreateCmd)
	InstanceCmd.AddCommand(instanceShowCmd)
	InstanceCmd.AddCommand(instanceUpdateCmd)
	InstanceCmd.AddCommand(instanceRemoveCmd)
	InstanceCmd.AddCommand(instanceRebootCmd)
	InstanceCmd.AddCommand(instanceSoftRebootCmd)
	// InstanceCmd.AddCommand(instanceConsoleCmd)
	InstanceCmd.AddCommand(instanceStopCmd)
	InstanceCmd.AddCommand(instanceStartCmd)
	InstanceCmd.AddCommand(instanceUpgradeCmd)
	// InstanceCmd.AddCommand(instanceMoveIPCmd)
	InstanceCmd.AddCommand(instanceSetFirewallCmd)
	InstanceCmd.AddCommand(instancePublicIPCmd)
	InstanceCmd.AddCommand(instancePasswordCmd)
	InstanceCmd.AddCommand(instanceTagCmd)
	InstanceCmd.AddCommand(instanceVncCmd)
	InstanceCmd.AddCommand(instanceRecoveryCmd)
	InstanceCmd.AddCommand(instanceRecoveryStatusCmd)
	InstanceCmd.AddCommand(snapshotCmd)
	InstanceCmd.AddCommand(instanceAllowedIPsUpdateCmd)
	InstanceCmd.AddCommand(instanceBandwidthUpdateCmd)


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
	instanceCreateCmd.Flags().StringVarP(&privateIPv4, "private_ipv4", "", "", "Private IPv4 address")
	instanceCreateCmd.Flags().StringVarP(&reservedIPv4, "reservedip", "", "", "Reserved IPv4 address")
	instanceCreateCmd.Flags().StringVar(&script, "script", "", "path to a script that will be uploaded to /usr/local/bin/civo-user-init-script on your instance, read/write/executable only by root and then will be executed at the end of the cloud initialization")
	instanceCreateCmd.Flags().BoolVar(&skipShebangCheck, "skip-shebang-check", false, "skip the shebang line check when passing a user init script")
	instanceCreateCmd.Flags().StringSliceVarP(&volumes, "volumes", "v", []string{}, "List of volumes to attach at boot")
	instanceCreateCmd.Flags().StringVarP(&volumetype, "volume-type", "", "", "Specify the volume type for the instance")
	instanceCreateCmd.Flags().StringArrayVar(&allowedIPs, "allowed-ips", []string{}, "A comma separated list of IP addresses that the instance is allowed to use")
	instanceCreateCmd.Flags().IntVar(&networkBandwidthLimit, "network-bandwidth-limit", 0, "The network bandwidth limit for the instance in Mbps (0 for unlimited)")

	instanceVncCmd.Flags().StringVarP(&duration, "duration", "d", "", "Duration for VNC access (e.g. 30m, 1h, 24h)")

	instanceStopCmd.Flags().BoolVarP(&waitStop, "wait", "w", false, "wait until the instance's is stoped")

	instanceAllowedIPsUpdateCmd.Flags().StringSliceVarP(&allowedIPsUpdate, "ips", "", []string{}, "Comma-separated list of IP addresses to allow (e.g., --ips 1.2.3.4,5.6.7.8). To clear all IPs, provide an empty string or omit the flag if the API interprets an empty list as clearing.")

	instanceBandwidthUpdateCmd.Flags().IntVarP(&bandwidthLimitUpdate, "limit", "l", 0, "Network bandwidth limit in Mbps (e.g., 1000). Use 0 for unlimited")
}
