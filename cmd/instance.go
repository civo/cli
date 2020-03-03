package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// instanceCmd represents the region command
var instanceCmd = &cobra.Command{
	Use:     "instance",
	Aliases: []string{"instances"},
	Short:   "Details of Civo instances",
}

// instanceListCmd represents the command to list current instances
var instanceListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List instances",
	Long: `List all current instances, including which is the default.
If you wish to use a custom format, the available fields are:

* ID
* OpenstackServerID
* Hostname
* ReverseDNS
* Size
* Region
* NetworkID
* PrivateIP
* PublicIP
* PseudoIP
* TemplateID
* SnapshotID
* InitialUser
* SSHKey
* Status
* Notes
* FirewallID
* Tags
* CivostatsdToken
* CivostatsdStats
* RescuePassword
* Script
* CreatedAt

Example: civo instance ls -o custom -f "ID: Name (PublicIP)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			return
		}

		instances, err := client.ListAllInstances()
		if err != nil {
			fmt.Printf("Unable to list instances: %s\n", aurora.Red(err))
			return
		}

		ow := utility.NewOutputWriter()
		for _, instance := range instances {
			ow.StartLine()

			if OutputFormat == "json" || OutputFormat == "custom" {
				ow.AppendData("ID", instance.ID)
				ow.AppendData("OpenstackServerID", instance.OpenstackServerID)
				ow.AppendData("Hostname", instance.Hostname)
				ow.AppendData("Size", instance.Size)
				ow.AppendData("Region", instance.Region)
				ow.AppendData("NetworkID", instance.NetworkID)
				ow.AppendData("PrivateIP", instance.PrivateIP)
				ow.AppendData("PublicIP", instance.PublicIP)
				ow.AppendData("PseudoIP", instance.PseudoIP)
				ow.AppendData("TemplateID", instance.TemplateID)
				ow.AppendData("SnapshotID", instance.SnapshotID)
				ow.AppendData("InitialUser", instance.InitialUser)
				ow.AppendData("SSHKey", instance.SSHKey)
				ow.AppendData("Status", instance.Status)
				ow.AppendData("Notes", instance.Notes)
				ow.AppendData("FirewallID", instance.FirewallID)
				ow.AppendData("Tags", strings.Join(instance.Tags, " "))
				ow.AppendData("CivostatsdToken", instance.CivostatsdToken)
				ow.AppendData("CivostatsdStats", instance.CivostatsdStats)
				ow.AppendData("Script", instance.Script)
				ow.AppendData("CreatedAt", instance.CreatedAt.Format(time.RFC1123))

				ow.AppendData("ReverseDNS", instance.ReverseDNS)
				ow.AppendData("PrivateIP", instance.PrivateIP)
				ow.AppendData("PublicIP", instance.PublicIP)
				ow.AppendData("PseudoIP", instance.PseudoIP)
			} else {
				ow.AppendData("ID", instance.ID)
				ow.AppendData("Hostname", instance.Hostname)
				ow.AppendData("Size", instance.Size)
				ow.AppendData("Region", instance.Region)
				ow.AppendData("Public IP", instance.PublicIP)
				ow.AppendData("Status", instance.Status)
			}
		}

		switch OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(OutputFields)
		default:
			ow.WriteTable()
		}
	},
}

// instanceShowCmd represents the command to show the details for an instance
var instanceShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Short:   "Show instance",
	Long: `Show your current instance.
If you wish to use a custom format, the available fields are:

* ID
* OpenstackServerID
* Hostname
* ReverseDNS
* Size
* Region
* NetworkID
* PrivateIP
* PublicIP
* PseudoIP
* TemplateID
* SnapshotID
* InitialUser
* InitialPassword
* SSHKey
* Status
* Notes
* FirewallID
* Tags
* CivostatsdToken
* CivostatsdStats
* RescuePassword
* Script
* CreatedAt

Example: civo instance show ID/NAME -o custom -f "Key1: Key2"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			return
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Unable to search instances: %s\n", aurora.Red(err))
			return
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		if OutputFormat == "json" || OutputFormat == "custom" {
			ow.AppendData("ID", instance.ID)
			ow.AppendData("OpenstackServerID", instance.OpenstackServerID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendData("Size", instance.Size)
			ow.AppendData("Region", instance.Region)
			ow.AppendData("NetworkID", instance.NetworkID)
			ow.AppendData("PrivateIP", instance.PrivateIP)
			ow.AppendData("PublicIP", instance.PublicIP)
			ow.AppendData("PseudoIP", instance.PseudoIP)
			ow.AppendData("TemplateID", instance.TemplateID)
			ow.AppendData("SnapshotID", instance.SnapshotID)
			ow.AppendData("InitialUser", instance.InitialUser)
			ow.AppendData("SSHKey", instance.SSHKey)
			ow.AppendData("Status", instance.Status)
			ow.AppendData("Notes", instance.Notes)
			ow.AppendData("FirewallID", instance.FirewallID)
			ow.AppendData("Tags", strings.Join(instance.Tags, " "))
			ow.AppendData("CivostatsdToken", instance.CivostatsdToken)
			ow.AppendData("CivostatsdStats", instance.CivostatsdStats)
			ow.AppendData("Script", instance.Script)
			ow.AppendData("CreatedAt", instance.CreatedAt.Format(time.RFC1123))

			ow.AppendData("ReverseDNS", instance.ReverseDNS)
			ow.AppendData("PrivateIP", instance.PrivateIP)
			ow.AppendData("PublicIP", instance.PublicIP)
			ow.AppendData("PseudoIP", instance.PseudoIP)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		} else {
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Openstack Server ID", instance.OpenstackServerID)
			ow.AppendData("Hostname", instance.Hostname)
			if instance.Hostname != instance.ReverseDNS && instance.ReverseDNS != "" {
				ow.AppendData("Reverse DNS", instance.ReverseDNS)
			}
			ow.AppendData("Size", instance.Size)
			ow.AppendData("Region", instance.Region)
			ow.AppendData("Network ID", instance.NetworkID)
			if instance.PseudoIP != "" {
				PublicIP := fmt.Sprintf("%s => %s", instance.PseudoIP, instance.PublicIP)
				ow.AppendData("Public IP", PublicIP)
			} else {
				ow.AppendData("Public IP", instance.PublicIP)
			}
			ow.AppendData("Template ID", instance.TemplateID)
			ow.AppendData("Snapshot ID", instance.SnapshotID)
			ow.AppendData("Initial User", instance.InitialUser)
			ow.AppendData("SSH Key", instance.SSHKey)
			ow.AppendData("Status", instance.Status)
			ow.AppendData("Firewall ID", instance.FirewallID)
			ow.AppendData("Tags", strings.Join(instance.Tags, " "))
			ow.AppendData("Created At", instance.CreatedAt.Format(time.RFC1123))

			ow.WriteKeyValues()

			if len(instance.Notes) > 0 {
				fmt.Println()
				ow.WriteSubheader("Notes")
				fmt.Println(instance.Notes)
			}

			if len(instance.Script) > 0 {
				fmt.Println()
				ow.WriteSubheader("Script")
				fmt.Println(instance.Script)
			}
		}
	},
}

// TODO: instance create [--name=HOSTNAME] [...] -- create a new instance with specified hostname and provided options
// TODO: instance tags ID/HOSTNAME 'tag1 tag2 tag3...' -- retag instance by ID (input no tags to clear all tags)
// TODO: instance update ID/HOSTNAME [--name] [--notes] -- update details of instance
// TODO: instance remove ID/HOSTNAME -- removes an instance with ID/hostname entered (use with caution!) [delete, destroy, rm]

// TODO: instance reboot ID/HOSTNAME -- reboots instance with ID/hostname entered [hard-reboot]
// TODO: instance soft-reboot ID/HOSTNAME -- soft-reboots instance with ID entered
// TODO: instance console ID/HOSTNAME -- outputs a URL for a web-based console for instance with ID provided
// TODO: instance stop ID/HOSTNAME -- shuts down the instance with ID provided
// TODO: instance start ID/HOSTNAME -- starts a stopped instance with ID provided
// TODO: instance upgrade ID/HOSTNAME new-size -- Upgrade instance with ID to size provided (see civo sizes for size names)
// TODO: instance move-ip ID/HOSTNAME IP_Address -- move a public IP_Address to target instance
// TODO: instance firewall ID/HOSTNAME firewall_id -- set instance with ID/HOSTNAME to use firewall with firewall_id
// TODO: instance public_ip ID/HOSTNAME -- Show public IP of ID/hostname [ip]
// TODO: instance password ID/HOSTNAME -- Show the default user password for instance with ID/HOSTNAME

func init() {
	rootCmd.AddCommand(instanceCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceShowCmd)
}
