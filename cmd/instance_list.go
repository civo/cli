package cmd

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceListCmd = &cobra.Command{
	Use:     "ls",
	Example: `civo instance ls -o custom -f "ID: Name (PublicIP)"`,
	Aliases: []string{"list", "all"},
	Short:   "List instances",
	Long: `List all current instances.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
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
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instances, err := client.ListAllInstances()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, instance := range instances {
			ow.StartLine()

			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendData("Region", client.Region)
			ow.AppendData("Size", instance.Size)
			ow.AppendDataWithLabel("CPUCores", strconv.Itoa(instance.CPUCores), "Cpu Cores")
			ow.AppendDataWithLabel("RAMMegabytes", strconv.Itoa(instance.RAMMegabytes), "Ram")
			ow.AppendDataWithLabel("DiskGigabytes", strconv.Itoa(instance.DiskGigabytes), "SSD disk")
			ow.AppendDataWithLabel("PublicIP", instance.PublicIP, "Public IP")
			ow.AppendDataWithLabel("PrivateIP", instance.PrivateIP, "Private IP")
			ow.AppendData("Status", utility.ColorStatus(instance.Status))

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendData("Status", instance.Status)
				ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
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

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
