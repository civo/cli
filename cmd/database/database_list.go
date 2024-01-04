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

var dbListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo db ls`,
	Short:   "List all databases",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		databases, err := client.ListDatabases()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, db := range databases.Items {
			ports := []string{}
			for _, user := range db.DatabaseUserInfo {
				ports = append(ports, fmt.Sprintf("%d", user.Port))
			}

			ow.StartLine()
			ow.AppendDataWithLabel("id", utility.TrimID(db.ID), "ID")
			ow.AppendDataWithLabel("name", db.Name, "Name")
			ow.AppendDataWithLabel("size", db.Size, "Size")
			ow.AppendDataWithLabel("nodes", strconv.Itoa(db.Nodes), "Nodes")
			ow.AppendDataWithLabel("software", db.Software, "Software")
			ow.AppendDataWithLabel("software_version", db.SoftwareVersion, "Software Version")
			ow.AppendDataWithLabel("host", db.PublicIPv4, "Host")
			ow.AppendDataWithLabel("port", strings.Join(ports, ","), "Port")
			ow.AppendDataWithLabel("status", db.Status, "Status")

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("firewall_id", db.FirewallID, "Firewall ID")
				ow.AppendDataWithLabel("network_id", db.NetworkID, "Network ID")
				ow.AppendDataWithLabel("id", db.ID, "ID")
			}
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
			fmt.Println("To get the credentials for a database, use `civo db credential <name/ID>`")
		}
	},
}
