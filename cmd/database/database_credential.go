package database

import (
	"fmt"
	"os"
	"strings"

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
			for _, userInfo := range db.DatabaseUserInfo {

				ow.StartLine()
				fmt.Println()
				ow.AppendDataWithLabel("database_id", utility.TrimID(db.ID), "Database ID")
				ow.AppendDataWithLabel("name", db.Name, "Name")
				ow.AppendDataWithLabel("host", db.PublicIPv4, "Host")
				ow.AppendDataWithLabel("port", fmt.Sprintf("%d", userInfo.Port), "Port")
				ow.AppendDataWithLabel("username", userInfo.Username, "Username")
				ow.AppendDataWithLabel("password", userInfo.Password, "Password")

				if common.OutputFormat == "json" || common.OutputFormat == "custom" {
					ow.AppendDataWithLabel("firewall_id", db.FirewallID, "Firewall ID")
					ow.AppendDataWithLabel("network_id", db.NetworkID, "Network ID")
					ow.AppendDataWithLabel("database_id", db.ID, "Database ID")

				}
			}
		} else {
			for _, user := range db.DatabaseUserInfo {
				fmt.Printf("%s://%s:%s@%s:%d\n", strings.ToLower(db.Software), user.Username, user.Password, db.PublicIPv4, user.Port)
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
