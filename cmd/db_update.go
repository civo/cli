package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"change", "modify"},
	Example: "civo db update OLD_NAME NEW_NAME",
	Short:   "Rename a database",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		oldDB, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("Database %s", err)
			os.Exit(1)
		}

		updateDB := civogo.UpdateDatabaseRequest{
			Name: args[1],
			ID:   oldDB.ID,
		}

		db, err := client.UpdateDatabase(&updateDB)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": db.ID, "name": db.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Renamed the database called %s with ID %s to %s\n", utility.Green(oldDB.Name), utility.Green(db.ID), utility.Green(db.Name))
		}
	},
}
