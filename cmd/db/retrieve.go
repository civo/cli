package db

import (
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var getCommand = &cobra.Command{
	Use:     "get",
	Short:   "Get details about a Civo Database",
	Aliases: []string{"show", "inspect"},
	Example: `civo db show ID/HOSTNAME"`,
	Args:    cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllDatabasesList(), cobra.ShellCompDirectiveNoFileComp
		}
		return getDatabasesWithSimilarParam(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		reg, err := utility.GetCurrentRegion()
		if err != nil {
			utility.Error("Failed to get region: %s", err)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if reg != "" {
			client.Region = reg
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		db, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("Database %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", db.ID)
		ow.AppendData("Name", db.Name)
		ow.AppendData("Region", client.Region)
		ow.AppendData("Replicas", strconv.Itoa(db.Replicas))
		ow.AppendData("Size", db.Size)
		ow.AppendData("Status", db.Status)

		ow.AppendDataWithLabel("Public IP", db.PublicIP, "Public IP")
		ow.WriteKeyValues()
		// TODO: Figure out JSON formatting, etc
	},
}

func getAllDatabasesList() []string {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	database, err := client.ListDatabases()
	if err != nil {
		utility.Error("Unable to list databases %s", err)
		os.Exit(1)
	}

	var databaseList []string
	for _, v := range database.Items {
		databaseList = append(databaseList, v.Name)
	}

	return databaseList
}

// getAllDatabasesList returns a list of all the databases with given string in their name or ID
func getDatabasesWithSimilarParam(value string) []string {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	database, err := client.FindDatabase(value)
	if err != nil {
		utility.Error("Unable to list databases %s", err)
		os.Exit(1)
	}

	var dbList []string
	dbList = append(dbList, database.Name)

	return dbList

}
