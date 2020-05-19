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
	Example: "civo instance create --hostname=foo.example.com [flags]",
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
	* CreatedAt`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		config, err := client.NewInstanceConfig()
		if err != nil {
			utility.Error("Unable to create a new config for the instance %s", err)
			os.Exit(1)
		}

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
			s.Prefix = fmt.Sprintf("Creating instance (%s)... ", resp.Hostname)
			s.Start()

			for stillCreating {
				instanceCheck, err := client.FindInstance(resp.ID)
				if err != nil {
					utility.Error("Unable to find the network %s", err)
					os.Exit(1)
				}
				if instanceCheck.Status == "ACTIVE" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
		}

		// we look for the created instance to obtain the data that we need
		// like PublicIP
		instance, err := client.FindInstance(resp.ID)
		if err != nil {
			utility.Error("Unable to find the instance %s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) has been created\n", utility.Green(instance.Hostname), instance.PublicIP)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", resp.ID)
			ow.AppendData("Hostname", resp.Hostname)
			ow.AppendData("Size", resp.Size)
			ow.AppendData("Region", resp.Region)
			ow.AppendDataWithLabel("PublicIP", resp.PublicIP, "Public IP")
			ow.AppendData("Status", resp.Status)
			ow.AppendDataWithLabel("OpenstackServerID", resp.OpenstackServerID, "Openstack Server ID")
			ow.AppendData("NetworkID", resp.NetworkID)
			ow.AppendData("TemplateID", resp.TemplateID)
			ow.AppendData("SnapshotID", resp.SnapshotID)
			ow.AppendData("InitialUser", resp.InitialUser)
			ow.AppendData("SSHKey", resp.SSHKey)
			ow.AppendData("Notes", resp.Notes)
			ow.AppendData("FirewallID", resp.FirewallID)
			ow.AppendData("Tags", strings.Join(resp.Tags, " "))
			ow.AppendData("CivostatsdToken", resp.CivostatsdToken)
			ow.AppendData("CivostatsdStats", resp.CivostatsdStats)
			ow.AppendData("Script", resp.Script)
			ow.AppendData("CreatedAt", resp.CreatedAt.Format(time.RFC1123))
			ow.AppendData("ReverseDNS", resp.ReverseDNS)
			ow.AppendData("PrivateIP", resp.PrivateIP)
			ow.AppendData("PublicIP", resp.PublicIP)
			ow.AppendData("PseudoIP", resp.PseudoIP)

			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
