package cmd

import (
	"fmt"
	"os"
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
			os.Exit(1)
		}

		instances, err := client.ListAllInstances()
		if err != nil {
			fmt.Printf("Unable to list instances: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, instance := range instances {
			ow.StartLine()

			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendData("Size", instance.Size)
			ow.AppendData("Region", instance.Region)
			ow.AppendDataWithLabel("PublicIP", instance.PublicIP, "Public IP")
			ow.AppendData("Status", instance.Status)

			if OutputFormat == "json" || OutputFormat == "custom" {
				ow.AppendData("OpenstackServerID", instance.OpenstackServerID)
				ow.AppendData("NetworkID", instance.NetworkID)
				ow.AppendData("PrivateIP", instance.PrivateIP)
				ow.AppendData("PublicIP", instance.PublicIP)
				ow.AppendData("TemplateID", instance.TemplateID)
				ow.AppendData("SnapshotID", instance.SnapshotID)
				ow.AppendData("InitialUser", instance.InitialUser)
				ow.AppendData("SSHKey", instance.SSHKey)
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
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Unable to search instances: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", instance.ID)
		ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
		ow.AppendData("Hostname", instance.Hostname)
		ow.AppendData("Status", instance.Status)
		ow.AppendData("Size", instance.Size)
		ow.AppendData("Region", instance.Region)
		ow.AppendDataWithLabel("NetworkID", instance.NetworkID, "Network ID")
		ow.AppendDataWithLabel("TemplateID", instance.TemplateID, "Template ID")
		ow.AppendDataWithLabel("SnapshotID", instance.SnapshotID, "Snapshot ID")
		ow.AppendDataWithLabel("InitialUser", instance.InitialUser, "Initial User")
		ow.AppendDataWithLabel("SSHKey", instance.SSHKey, "SSH Key")
		ow.AppendDataWithLabel("FirewallID", instance.FirewallID, "Firewall ID")
		ow.AppendData("Tags", strings.Join(instance.Tags, " "))
		ow.AppendDataWithLabel("CreatedAt", instance.CreatedAt.Format(time.RFC1123), "Created At")
		ow.AppendDataWithLabel("PrivateIP", instance.PrivateIP, "Private IP")

		if OutputFormat == "json" || OutputFormat == "custom" {
			ow.AppendDataWithLabel("PublicIP", instance.PublicIP, "Public IP")
			ow.AppendDataWithLabel("PseudoIP", instance.PseudoIP, "Pseudo IP")
			ow.AppendData("Notes", instance.Notes)
			ow.AppendData("Script", instance.Script)

			ow.AppendData("ReverseDNS", instance.ReverseDNS)
			ow.AppendData("PublicIP", instance.PublicIP)
			ow.AppendData("PseudoIP", instance.PseudoIP)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		} else {
			if instance.PseudoIP != "" {
				PublicIP := fmt.Sprintf("%s => %s", instance.PseudoIP, instance.PublicIP)
				ow.AppendData("Public IP", PublicIP)
			} else {
				ow.AppendData("Public IP", instance.PublicIP)
			}
			if instance.Hostname != instance.ReverseDNS && instance.ReverseDNS != "" {
				ow.AppendData("Reverse DNS", instance.ReverseDNS)
			}

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

var instanceRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "destroy", "rm"},
	Short:   "Remove/delete instance",
	Long: `Remove the specified instance by part of the ID or name.

Example: civo instance remove ID/NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.DeleteInstance(instance.ID)
		if err != nil {
			fmt.Printf("Removing instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) has been removed\n", aurora.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
			ow.AppendData("Hostname", instance.Hostname)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}

var instanceRebootCmd = &cobra.Command{
	Use:     "reboot",
	Aliases: []string{"hard-reboot"},
	Short:   "Hard reboot an instance",
	Long: `Pull the power and restart the specified instance by part of the ID or name.

Example: civo instance reboot ID/NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.RebootInstance(instance.ID)
		if err != nil {
			fmt.Printf("Rebooting instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being rebooted\n", aurora.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
			ow.AppendData("Hostname", instance.Hostname)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}

// TODO: instance soft-reboot ID/HOSTNAME -- soft-reboots instance with ID entered
var instanceSoftRebootCmd = &cobra.Command{
	Use:   "soft-reboot",
	Short: "Hard reboot an instance",
	Long: `Pull the power and restart the specified instance by part of the ID or name.

Example: civo instance reboot ID/NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.SoftRebootInstance(instance.ID)
		if err != nil {
			fmt.Printf("Soft-rebooting instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being soft-rebooted\n", aurora.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
			ow.AppendData("Hostname", instance.Hostname)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}

// TODO: instance console ID/HOSTNAME -- outputs a URL for a web-based console for instance with ID provided
var instanceConsoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Get console URL for instance",
	Long: `Get the web console's URL for a given instance by part of the ID or name.

Example: civo instance console ID/NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		url, err := client.GetInstanceConsoleURL(instance.ID)
		if err != nil {
			fmt.Printf("Getting console URL: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) has a console at %s\n", aurora.Green(instance.Hostname), instance.ID,
				aurora.Green(url))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendDataWithLabel("URL", url, "Console URL")
			ow.AppendData("Hostname", instance.Hostname)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}

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
	instanceCmd.AddCommand(instanceRemoveCmd)
	instanceCmd.AddCommand(instanceRebootCmd)
	instanceCmd.AddCommand(instanceSoftRebootCmd)
	instanceCmd.AddCommand(instanceConsoleCmd)
}
