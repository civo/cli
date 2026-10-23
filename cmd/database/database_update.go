package database

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"modify", "change"},
	Short:   "Update a database",
	Example: `civo db update DB_NAME --nodes 5 --name NEW_NAME
civo db update DB_NAME --storage 120
civo db update DB_NAME --extensions pg_cron,postgis
civo db update DB_NAME --parameter max_connections=200 --parameter shared_buffers=1GB`,
	Args: cobra.MinimumNArgs(1),
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

		findDB, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("Database %s", err)
			os.Exit(1)
		}

		updatedDB := findDB
		if cmd.Flags().Changed("name") || cmd.Flags().Changed("nodes") || cmd.Flags().Changed("firewall") {
			updatedDB, err = client.UpdateDatabase(findDB.ID, &civogo.UpdateDatabaseRequest{
				Name:       updatedName,
				Nodes:      &nodes,
				FirewallID: firewallID,
				Region:     client.Region,
			})
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		// Storage, extensions and parameters live on the v3alpha1 API and go
		// out in one request, so a failure cannot leave them half applied.
		v3 := &civogo.UpdateDatabaseV3alpha1Request{Region: client.Region}
		v3Changed := false
		if storageGB > 0 {
			gb := storageGB
			v3.StorageGB = &gb
			v3Changed = true
		}
		if cmd.Flags().Changed("extensions") {
			update := true
			v3.Extensions = parseExtensions(extensions)
			v3.UpdateExtensions = &update
			v3Changed = true
		}
		if cmd.Flags().Changed("parameter") {
			params, err := parseDatabaseParameters(parameters)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			update := true
			v3.DatabaseParameters = params
			v3.UpdateDatabaseParameters = &update
			v3Changed = true
		}
		if v3Changed {
			updatedDB, err = client.UpdateDatabaseV3alpha1(findDB.ID, v3)
			if err != nil {
				utility.Error("Failed to update the database: %s", err)
				os.Exit(1)
			}
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": updatedDB.ID, "name": updatedDB.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The Database %s was updated\n", utility.Green(findDB.Name))
			showDatabaseDeprecationWarnings(*findDB)
		}
	},
}
