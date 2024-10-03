package instance

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/civo/cli/common"
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

	* id
	* hostname
	* status
	* size
	* cpu_cores
	* ram_mb
	* disk_gb
	* region
	* network_id
	* diskimage_id
	* initial_user
	* initial_password
	* ssh_key
	* firewall_id
	* tags
	* created_at
	* private_ip
	* public_ip
	* notes
	* script
	* reverse_dns`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
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

		ow.AppendDataWithLabel("id", instance.ID, "ID")
		ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
		ow.AppendDataWithLabel("status", utility.ColorStatus(instance.Status), "Status")
		ow.AppendDataWithLabel("size", instance.Size, "Size")
		ow.AppendDataWithLabel("volume-type", instance.VolumeType, "Volume Type")
		ow.AppendDataWithLabel("cpu_cores", strconv.Itoa(instance.CPUCores), "Cpu Cores")
		ow.AppendDataWithLabel("ram_mb", strconv.Itoa(instance.RAMMegabytes), "Ram")
		ow.AppendDataWithLabel("disk_gb", strconv.Itoa(instance.DiskGigabytes), "SSD disk")
		ow.AppendDataWithLabel("region", client.Region, "Region")
		ow.AppendDataWithLabel("network_id", instance.NetworkID, "Network ID")
		ow.AppendDataWithLabel("diskimage_id", instance.SourceID, "Disk image ID")
		ow.AppendDataWithLabel("initial_user", instance.InitialUser, "Initial User")
		ow.AppendDataWithLabel("initial_password", instance.InitialPassword, "Initial Password")
		ow.AppendDataWithLabel("ssh_key_id", instance.SSHKeyID, "SSH Key ID")
		ow.AppendDataWithLabel("firewall_id", instance.FirewallID, "Firewall ID")
		ow.AppendDataWithLabel("tags", strings.Join(instance.Tags, " "), "Tags")
		ow.AppendDataWithLabel("created_at", instance.CreatedAt.Format(time.RFC1123), "Created At")
		ow.AppendDataWithLabel("private_ip", instance.PrivateIP, "Private IP")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("public_ip", instance.PublicIP, "Public IP")
			ow.AppendDataWithLabel("notes", instance.Notes, "notes")
			ow.AppendDataWithLabel("script", instance.Script, "Script")

			ow.AppendDataWithLabel("reverse_dns", instance.ReverseDNS, "Reverse DNS")
			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		} else {
			if instance.PseudoIP != "" {
				publicIP := fmt.Sprintf("%s => %s", instance.PrivateIP, instance.PublicIP)
				ow.AppendDataWithLabel("public_ip", publicIP, "Public IP")
			} else {
				ow.AppendDataWithLabel("public_ip", instance.PublicIP, "Public IP")
			}
			if instance.Hostname != instance.ReverseDNS && instance.ReverseDNS != "" {
				ow.AppendDataWithLabel("reverse_dns", instance.ReverseDNS, "Reverse DNS")
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
				sDec, err := b64.StdEncoding.DecodeString(instance.Script)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				fmt.Println(string(sDec))
			}
		}
	},
}
