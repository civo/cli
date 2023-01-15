// Package database is the root command for civo database.
package database

import (
	"errors"

	"github.com/spf13/cobra"
)

var firewallID, networkID, size, updatedName string
var nodes int

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

	dbCredentialCmd.Flags().BoolVarP(&connectionString, "connection-string", "c", false, "show the connection string for the database")

	dbCreateCmd.Flags().IntVarP(&nodes, "nodes", "", 1, "the number of nodes for the database")
	dbCreateCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the database")
	dbCreateCmd.Flags().StringVarP(&networkID, "network", "n", "", "the network to use for the database")
	dbCreateCmd.Flags().StringVarP(&rulesFirewall, "firewall-rules", "u", "", "the firewall rules to use for the database")
	dbCreateCmd.Flags().StringVarP(&size, "size", "s", "g3.db.small", "the size of the database. You can list available DB sizes by `civo size list -s database`")

	dbUpdateCmd.Flags().IntVarP(&nodes, "nodes", "", 0, "the number of nodes for the database")
	dbUpdateCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the database")
	dbUpdateCmd.Flags().StringVarP(&updatedName, "name", "n", "", "the new name for the database")
}
