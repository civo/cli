package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Short:   "Show instance",
	Long: `Show your current instance.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* OpenstackServerID
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
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Unable to search instances %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", instance.ID)
		ow.AppendData("Hostname", instance.Hostname)
		ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
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
