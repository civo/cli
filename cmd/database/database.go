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
	DBCmd.AddCommand(dbForkCmd)
	DBCmd.AddCommand(dbExtensionCmd)
	DBCmd.AddCommand(dbUpgradeCmd)

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
	dbUpdateCmd.Flags().IntVarP(&storageGB, "storage", "", 0, "grow the storage of every node to this many GB (storage can only be increased; GB beyond the plan are billed as an add-on)")
	dbUpdateCmd.Flags().StringVarP(&extensions, "extensions", "", "", "comma separated list of PostgreSQL extensions to enable, replacing the current set (an empty value disables them all); see civo db extension list")
	dbUpdateCmd.Flags().StringArrayVarP(&parameters, "parameter", "p", nil, "database parameter as key=value, repeatable; replaces the current custom parameters")

	dbRestoreCmd.Flags().StringVarP(&backup, "backup", "b", "", "the backup name which you can restore database")
	dbRestoreCmd.Flags().StringVarP(&restoreName, "name", "n", "", "name of the restore")
	dbRestoreCmd.Flags().StringVarP(&targetTime, "target-time", "t", "", "restore to this point in time instead of a backup, as RFC 3339 (e.g. 2026-10-21T10:00:00Z); must not be earlier than the restore window shown by civo db show")
	dbRestoreCmd.MarkFlagsMutuallyExclusive("backup", "target-time")
	dbRestoreCmd.MarkFlagsOneRequired("backup", "target-time")
	_ = dbRestoreCmd.MarkFlagRequired("name")

	dbForkCmd.Flags().StringVarP(&forkName, "name", "n", "", "name of the new database")
	dbForkCmd.Flags().StringVarP(&targetTime, "target-time", "t", "", "fork the source as it was at this point in time, as RFC 3339 (e.g. 2026-10-21T10:00:00Z)")
	dbForkCmd.Flags().StringVarP(&forkSize, "size", "s", "", "size of the new database (defaults to the source's)")
	dbForkCmd.Flags().IntVarP(&forkNodes, "nodes", "", 0, "number of nodes of the new database (defaults to the source's)")
	dbForkCmd.Flags().StringVarP(&forkNetworkID, "network", "", "", "network of the new database (defaults to the source's)")
	dbForkCmd.Flags().StringVarP(&forkFirewallID, "firewall", "", "", "firewall of the new database")
	dbForkCmd.Flags().IntVarP(&storageGB, "storage", "", 0, "storage of every node of the new database in GB")
	dbForkCmd.Flags().StringVarP(&extensions, "extensions", "", "", "comma separated list of PostgreSQL extensions to enable on the new database")
	dbForkCmd.Flags().BoolVarP(&forkWait, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the new database to be Ready")
	_ = dbForkCmd.MarkFlagRequired("name")

	dbUpgradeCmd.Flags().StringVarP(&upgradeTargetVersion, "version", "v", "", "engine release to upgrade to (defaults to the newest available)")
	dbUpgradeCmd.Flags().BoolVar(&upgradeDryRun, "dry-run", false, "only show which databases would be upgraded and restarted")
}
