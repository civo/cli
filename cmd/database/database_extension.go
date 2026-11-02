package database

import (
	"errors"
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

// dbExtensionCmd is the root command for the db extension subcommand
var dbExtensionCmd = &cobra.Command{
	Use:     "extension",
	Aliases: []string{"extensions", "ext"},
	Short:   "Manage PostgreSQL extensions of Civo Databases",
	Long: `List the extensions a region offers. Enable them on a database with
"civo db update <DATABASE> --extensions pg_cron,postgis".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var dbExtensionListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "all"},
	Short:   "List the PostgreSQL extensions available in the current region",
	Example: "civo db extension list",
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

		exts, err := client.ListDatabaseExtensions()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, ext := range exts {
			ow.StartLine()
			ow.AppendDataWithLabel("name", ext.Name, "Name")
			ow.AppendDataWithLabel("description", ext.Description, "Description")
			ow.AppendDataWithLabel("versions", strings.Join(ext.Versions, ", "), "Versions")
			restart := "No"
			if ext.RequiresRestart {
				restart = "Yes"
			}
			ow.AppendDataWithLabel("requires_restart", restart, "Requires Restart")
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

func init() {
	dbExtensionCmd.AddCommand(dbExtensionListCmd)
}
