package database

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var connectionString bool

var dbCredentialCmd = &cobra.Command{
	Use:     "credential",
	Example: `civo db credential ID/NAME`,
	Aliases: []string{"credentials", "creds", "cred"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show credential details of a database",
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

		ow := utility.NewOutputWriter()

		if !connectionString {
			ow.StartLine()
			fmt.Println()
			ow.AppendDataWithLabel("id", db.ID, "ID")
			ow.AppendDataWithLabel("name", db.Name, "Name")
			ow.AppendDataWithLabel("host", fmt.Sprintf("%s:%d", db.PublicIPv4, db.Port), "Host")
			ow.AppendDataWithLabel("username", db.Username, "Username")
			ow.AppendDataWithLabel("password", db.Password, "Password")

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("firewall_id", db.FirewallID, "Firewall ID")
				ow.AppendDataWithLabel("network_id", db.NetworkID, "Network ID")
			}
		} else {
			fmt.Printf("mysql://%s:%s@%s:%d\n", db.Username, db.Password, db.PublicIPv4, db.Port)
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
