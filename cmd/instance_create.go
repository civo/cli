package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"strings"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var wait bool
var hostnameCreate, size, template, snapshot, publicip, initialuser, sshkey, tags, network string

var instanceCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a new instance",
	Long: `You can create an instance with a hostname parameter, as well as any options you provide.
If you wish to use a custom format, the available fields are:

	* ID
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

Example: civo instance create --hostname=foo.example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		config, err := client.NewInstanceConfig()

		if hostnameCreate != "" {
			config.Hostname = hostnameCreate
		}

		if size != "" {
			config.Size = size
		}

		if template != "" {
			config.TemplateID = template
		}

		if snapshot != "" {
			config.SnapshotID = snapshot
		}

		if publicip != "" {
			config.PublicIPRequired = publicip
		}

		if initialuser != "" {
			config.InitialUser = initialuser
		}

		if sshkey != "" {
			sshKey, err := client.FindSSHKey(sshkey)
			if err != nil {
				utility.Error("Unable to find the ssh key %s", err)
				os.Exit(1)
			}
			config.SSHKeyID = sshKey.ID
		}

		if network != "" {
			net, err := client.FindNetwork(network)
			if err != nil {
				utility.Error("Unable to find the network %s", err)
				os.Exit(1)
			}
			config.NetworkID = net.ID
		}

		if tags != "" {
			config.TagsList = tags
		}

		resp, err := client.CreateInstance(config)
		if err != nil {
			utility.Error("error creating instance %s", err)
			os.Exit(1)
		}

		if wait == true {

			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Creating instance... "
			s.Start()

			for stillCreating {
				instanceCheck, _ := client.FindInstance(resp.ID)
				if instanceCheck.Status == "ACTIVE" {
					stillCreating = false
					s.Stop()
				}
				time.Sleep(5 * time.Second)
			}
		}

		instance, _ := client.FindInstance(resp.ID)

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) has been create\n", utility.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendData("Size", instance.Size)
			ow.AppendData("Region", instance.Region)
			ow.AppendDataWithLabel("PublicIP", instance.PublicIP, "Public IP")
			ow.AppendData("Status", instance.Status)
			ow.AppendDataWithLabel("OpenstackServerID", instance.OpenstackServerID, "Openstack Server ID")
			ow.AppendData("NetworkID", instance.NetworkID)
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

			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
