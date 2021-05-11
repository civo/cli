package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceShowCmd = &cobra.Command{
	Use:     "show",
	Example: `civo instance show ID/HOSTNAME`,
	Aliases: []string{"get", "inspect"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show instance",
	Long: `View the details for an instance.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* OpenstackServerID
	* ReverseDNS
	* Size
	* CPUCores
	* RAMMegabytes
	* DiskGigabytes
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
	* CreatedAt`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", instance.ID)
		ow.AppendData("Hostname", instance.Hostname)
		ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
		ow.AppendData("Status", instance.Status)
		ow.AppendData("Size", instance.Size)
		ow.AppendDataWithLabel("CPUCores", strconv.Itoa(instance.CPUCores), "Cpu Cores")
		ow.AppendDataWithLabel("RAMMegabytes", strconv.Itoa(instance.RAMMegabytes), "Ram")
		ow.AppendDataWithLabel("DiskGigabytes", strconv.Itoa(instance.DiskGigabytes), "SSD disk")
		ow.AppendData("Region", client.Region)
		ow.AppendDataWithLabel("NetworkID", instance.NetworkID, "Network ID")
		ow.AppendDataWithLabel("TemplateID", instance.TemplateID, "Template ID")
		ow.AppendDataWithLabel("SnapshotID", instance.SnapshotID, "Snapshot ID")
		ow.AppendDataWithLabel("InitialUser", instance.InitialUser, "Initial User")
		ow.AppendDataWithLabel("SSHKey", instance.SSHKey, "SSH Key")
		ow.AppendDataWithLabel("FirewallID", instance.FirewallID, "Firewall ID")
		ow.AppendData("Tags", strings.Join(instance.Tags, " "))
		ow.AppendDataWithLabel("CreatedAt", instance.CreatedAt.Format(time.RFC1123), "Created At")
		ow.AppendDataWithLabel("PrivateIP", instance.PrivateIP, "Private IP")

		if outputFormat == "json" || outputFormat == "custom" {
			ow.AppendDataWithLabel("PublicIP", instance.PublicIP, "Public IP")
			ow.AppendDataWithLabel("PseudoIP", instance.PseudoIP, "Pseudo IP")
			ow.AppendData("Notes", instance.Notes)
			ow.AppendData("Script", instance.Script)

			ow.AppendData("ReverseDNS", instance.ReverseDNS)
			ow.AppendData("PublicIP", instance.PublicIP)
			ow.AppendData("PseudoIP", instance.PseudoIP)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		} else {
			if instance.PseudoIP != "" {
				publicIP := fmt.Sprintf("%s => %s", instance.PseudoIP, instance.PublicIP)
				ow.AppendData("Public IP", publicIP)
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
