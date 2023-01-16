package database

import (
	"fmt"
	"os"
	"strconv"

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
			ow.StartLine()
			ow.AppendDataWithLabel("id", utility.TrimID(db.ID), "ID")
			ow.AppendDataWithLabel("name", db.Name, "Name")
			ow.AppendDataWithLabel("size", db.Size, "Size")
			ow.AppendDataWithLabel("nodes", strconv.Itoa(db.Nodes), "Nodes")
			ow.AppendDataWithLabel("software", db.Software, "Software")
			ow.AppendDataWithLabel("software_version", db.SoftwareVersion, "Software Version")
			ow.AppendDataWithLabel("host", fmt.Sprintf("%s:%d", db.PublicIPv4, db.Port), "Host")
			ow.AppendDataWithLabel("status", db.Status, "Status")

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("username", db.Username, "Username")
				ow.AppendDataWithLabel("password", db.Password, "Password")
				ow.AppendDataWithLabel("firewall_id", db.FirewallID, "Firewall ID")
				ow.AppendDataWithLabel("network_id", db.NetworkID, "Network ID")
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
