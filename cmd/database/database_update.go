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
	Example: "civo db update DB_NAME --nodes 5 --name NEW_NAME",
	Args:    cobra.MinimumNArgs(1),
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

		updatedDB, err := client.UpdateDatabase(findDB.ID, &civogo.UpdateDatabaseRequest{
			Name:       updatedName,
			Nodes:      &nodes,
			FirewallID: firewallID,
			Region:     client.Region,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": updatedDB.ID, "name": updatedDB.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The Database %s was updated\n", utility.Green(findDB.Name))
			os.Exit(0)
		}
	},
}
