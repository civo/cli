package database

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbShowCmd = &cobra.Command{
	Use:     "show",
	Example: `civo db show ID/NAME`,
	Aliases: []string{"get", "inspect"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show details of a database",
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

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		db, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ports := []string{}
		for _, user := range db.DatabaseUserInfo {
			ports = append(ports, fmt.Sprintf("%d", user.Port))
		}

		ow := utility.NewOutputWriter()

		ow.StartLine()
		fmt.Println()
		ow.AppendDataWithLabel("id", db.ID, "ID")
		ow.AppendDataWithLabel("name", db.Name, "Name")
		ow.AppendDataWithLabel("status", db.Status, "Status")
		ow.AppendDataWithLabel("dns_entry", db.DNSEntry, "DNS")
		ow.AppendDataWithLabel("size", db.Size, "Size")
		ow.AppendDataWithLabel("nodes", strconv.Itoa(db.Nodes), "Nodes")
		ow.AppendDataWithLabel("software", db.Software, "Software")
		ow.AppendDataWithLabel("software_version", db.SoftwareVersion, "Software Version")
		ow.AppendDataWithLabel("host", db.PublicIPv4, "Host")
		ow.AppendDataWithLabel("port", strings.Join(ports, ","), "Port")
		ow.AppendDataWithLabel("private_ipv4", db.PrivateIPv4, "Private IPv4")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("firewall_id", db.FirewallID, "Firewall ID")
			ow.AppendDataWithLabel("network_id", db.NetworkID, "Network ID")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
			fmt.Println("To get the credentials, run : civo db credential", db.Name)
		}
	},
}
