// Package database is the root command for civo database.
package database

import (
	"errors"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var firewallID, networkID, size, updatedName, software, softwareVersion string
var nodes int

// showDatabaseDeprecationWarning displays a warning message if the database version is deprecated
func showDatabaseDeprecationWarnings(databases ...civogo.Database) {
	// We want to show one warning per version at max.
	pgWarning := false
	mysqlWarning := false

	for _, db := range databases {
		software := strings.ToLower(db.Software)

		if software == "mysql" && !mysqlWarning {
			utility.Warning("MySQL databases are deprecated and will be removed in a future release. Please consider checking the documentation https://www.civo.com/docs/database/mysql/dump-mysql to understand how to keep using MySQL with Civo")
			mysqlWarning = true
		}

		if software == "postgresql" && strings.HasPrefix(db.SoftwareVersion, "14") && !pgWarning {
			utility.Warning("PostgreSQL 14 is deprecated and will be removed in a future release. Please migrate to PostgreSQL 17. For migration guidance, see: https://www.civo.com/docs/database/postgresql/migrate-from-14-to-17")
			pgWarning = true
		}
	}
}

// DBCmd is the root command for the db subcommand
var DBCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db", "databases"},
	Short:   "Manage Civo Database ",
	Long:    `Create, update, delete, and list Civo Databases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	DBCmd.AddCommand(dbListCmd)
	DBCmd.AddCommand(dbCreateCmd)
	DBCmd.AddCommand(dbUpdateCmd)
	DBCmd.AddCommand(dbShowCmd)
	DBCmd.AddCommand(dbDeleteCmd)
	DBCmd.AddCommand(dbCredentialCmd)
	DBCmd.AddCommand(dbSizeCmd)
	DBCmd.AddCommand(dbEngineCmd)
	DBCmd.AddCommand(dbBackupCmd)
	DBCmd.AddCommand(dbRestoreCmd)
	DBCmd.AddCommand(dbVersionListCmd)

	dbCredentialCmd.Flags().BoolVarP(&connectionString, "connection-string", "c", false, "show the connection string for the database")

	dbCreateCmd.Flags().IntVarP(&nodes, "nodes", "", 1, "the number of nodes for the database")
	dbCreateCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the database")
	dbCreateCmd.Flags().StringVarP(&networkID, "network", "n", "", "the network to use for the database")
	dbCreateCmd.Flags().StringVarP(&rulesFirewall, "firewall-rules", "u", "", "the firewall rules to use for the database")
	dbCreateCmd.Flags().StringVarP(&size, "size", "s", "g3.db.small", "the size of the database. You can list available DB sizes by `civo size list -s database`")
	dbCreateCmd.Flags().StringVarP(&software, "software", "m", "PostgreSQL", "the software to use for the database.")
	dbCreateCmd.Flags().StringVarP(&softwareVersion, "version", "v", "", "the version of the software to use for the database.")
	dbCreateCmd.Flags().BoolVarP(&waitDatabase, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the database to be ACTIVE")

	dbUpdateCmd.Flags().IntVarP(&nodes, "nodes", "", 0, "the number of nodes for the database")
	dbUpdateCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the database")
	dbUpdateCmd.Flags().StringVarP(&updatedName, "name", "n", "", "the new name for the database")

	dbRestoreCmd.Flags().StringVarP(&backup, "backup", "b", "", "the backup name which you can restore database")
	dbRestoreCmd.Flags().StringVarP(&restoreName, "name", "n", "", "name of the restore")
	_ = dbRestoreCmd.MarkFlagRequired("backup")
	_ = dbRestoreCmd.MarkFlagRequired("name")
}
