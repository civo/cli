package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"

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

		// Add check for database status
		if db.Status == "Pending" {
			utility.Printf("The DB %s is currently being provisioned, please wait...\n", utility.Green(db.Name))

			var s *spinner.Spinner
			if !common.Quiet {
				s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
				s.Writer = os.Stderr
				s.Prefix = fmt.Sprintf("Waiting for database (%s)... ", db.Name)
				s.Start()
			}

			for db.Status == "Pending" {
				db, err = client.FindDatabase(args[0])
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				time.Sleep(2 * time.Second)
			}
			if !common.Quiet {
				s.Stop()
			}
		}

		connStr := strings.ToLower(db.Software) + "://" + db.DatabaseUserInfo[0].Username + ":" + db.DatabaseUserInfo[0].Password + "@" + db.PublicIPv4 + ":" + fmt.Sprintf("%d", db.DatabaseUserInfo[0].Port)

		if connectionString {
			for _, user := range db.DatabaseUserInfo {
				utility.Printf("%s://%s:%s@%s:%d\n", strings.ToLower(db.Software), user.Username, user.Password, db.PublicIPv4, user.Port)
			}
			return
		}
		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendDataWithLabel("id", utility.TrimID(db.ID), "ID")
		ow.AppendDataWithLabel("name", db.Name, "Name")
		ow.AppendDataWithLabel("host", db.PublicIPv4, "Host")
		ow.AppendDataWithLabel("connection-string", connStr, "Connection String")

		// Show credentials for each user
		for i, userInfo := range db.DatabaseUserInfo {
			if i == 0 {
				// For the first user, show credentials directly
				ow.AppendDataWithLabel("port", fmt.Sprintf("%d", userInfo.Port), "Port")
				ow.AppendDataWithLabel("username", userInfo.Username, "Username")
				ow.AppendDataWithLabel("password", userInfo.Password, "Password")
			} else {
				// For additional users, add a numbered suffix
				ow.AppendDataWithLabel(fmt.Sprintf("port_%d", i+1), fmt.Sprintf("%d", userInfo.Port), fmt.Sprintf("Port %d", i+1))
				ow.AppendDataWithLabel(fmt.Sprintf("username_%d", i+1), userInfo.Username, fmt.Sprintf("Username %d", i+1))
				ow.AppendDataWithLabel(fmt.Sprintf("password_%d", i+1), userInfo.Password, fmt.Sprintf("Password %d", i+1))
			}
		}

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("firewall_id", db.FirewallID, "Firewall ID")
			ow.AppendDataWithLabel("network_id", db.NetworkID, "Network ID")
			ow.AppendDataWithLabel("database_id", db.ID, "Database ID")

			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		} else {
			ow.WriteKeyValues()
		}
	},
}
