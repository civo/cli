package instance

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/civo/cli/common"
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

	* id
	* hostname
	* region
	* size
	* cpu_cores
	* ram_mb
	* disk_gb
	* public_ip
	* private_ip
	* status
	* network_id
	* diskimage_id
	* initial_user
	* ssh_key
	* notes
	* firewall_id
	* tags
	* script
	* created_at
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

		instances, err := client.ListAllInstances()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, instance := range instances {
			ow.StartLine()

			ow.AppendDataWithLabel("id", utility.TruncateID(instance.ID), "ID")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			ow.AppendDataWithLabel("region", client.Region, "Region")
			ow.AppendDataWithLabel("size", instance.Size, "Size")
			ow.AppendDataWithLabel("public_ip", instance.PublicIP, "Public IP")
			ow.AppendDataWithLabel("private_ip", instance.PrivateIP, "Private IP")
			ow.AppendDataWithLabel("ipv6", instance.IPv6, "IPv6")
			ow.AppendDataWithLabel("status", utility.ColorStatus(instance.Status), "Status")

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("network_id", instance.NetworkID, "Network ID")
				ow.AppendDataWithLabel("cpu_cores", strconv.Itoa(instance.CPUCores), "Cpu Cores")
				ow.AppendDataWithLabel("ram_mb", strconv.Itoa(instance.RAMMegabytes), "Ram")
				ow.AppendDataWithLabel("disk_gb", strconv.Itoa(instance.DiskGigabytes), "SSD disk")
				ow.AppendDataWithLabel("diskimage_id", instance.SourceID, "Disk image ID")
				ow.AppendDataWithLabel("initial_user", instance.InitialUser, "Initial User")
				ow.AppendDataWithLabel("ssh_key", instance.SSHKey, "SSH Key")
				ow.AppendDataWithLabel("notes", instance.Notes, "Notes")
				ow.AppendDataWithLabel("firewall_id", instance.FirewallID, "Firewall ID")
				ow.AppendDataWithLabel("tags", strings.Join(instance.Tags, " "), "Tags")
				ow.AppendDataWithLabel("script", instance.Script, "Script")
				ow.AppendDataWithLabel("created_at", instance.CreatedAt.Format(time.RFC1123), "Created At")
				ow.AppendDataWithLabel("status", instance.Status, "Status")
				ow.AppendDataWithLabel("reverse_dns", instance.ReverseDNS, "Reverse DNS")
			}
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
