package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/civo/civogo"
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
		utility.EnsureCurrentRegion()

		check, region, err := utility.CheckAvailability("iaas", regionSet)
		if err != nil {
			utility.Error("Error checking availability %s", err)
			os.Exit(1)
		}

		if check {
			utility.Error("Sorry you can't create a instance in the %s region", region)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
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

		if region != "" {
			config.Region = region
		}

		if size != "" {
			config.Size = size
		}

		if template != "" {
			templateID := ""
			if client.Region == "SVG1" {
				findTemplate, err := client.FindTemplate(template)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				templateID = findTemplate.ID
			} else {
				findTemplate, err := client.FindDiskImage(template)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				templateID = findTemplate.ID
			}

			config.TemplateID = templateID
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
				utility.Error("%s", err)
				os.Exit(1)
			}
			config.SSHKeyID = sshKey.ID
		}

		if network != "" {
			net, err := client.FindNetwork(network)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			config.NetworkID = net.ID
		}

		if tags != "" {
			config.TagsList = tags
		}

		var executionTime, publicIP string
		startTime := utility.StartTime()

		var instance *civogo.Instance
		resp, err := client.CreateInstance(config)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if wait {
			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating instance (%s)... ", resp.Hostname)
			s.Start()

			for stillCreating {
				instance, err = client.FindInstance(resp.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				if instance.Status == "ACTIVE" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
			publicIP = fmt.Sprintf("(%s)", instance.PublicIP)
			executionTime = utility.TrackTime(startTime)
		} else {
			// we look for the created instance to obtain the data that we need
			// like PublicIP
			instance, err = client.FindInstance(resp.ID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		if outputFormat == "human" {
			if executionTime != "" {
				fmt.Printf("The instance %s %s has been created in %s\n", utility.Green(instance.Hostname), publicIP, executionTime)
			} else {
				fmt.Printf("The instance %s has been created\n", utility.Green(instance.Hostname))
			}
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
