package database

import (
	"errors"

	"github.com/spf13/cobra"
)

var firewallID, networkID string
var replicas int

// DBCmd is the root command for the db subcommand
var DBCmd = &cobra.Command{
	Use:     "db",
	Aliases: []string{"database"},
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
	DBCmd.AddCommand(dbDeleteCmd)

	dbCreateCmd.Flags().IntVarP(&replicas, "replicas", "r", 0, "the number of replicas for the database")
	dbCreateCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the database")
	dbCreateCmd.Flags().StringVarP(&networkID, "network", "n", "", "the network to use for the database")

	dbUpdateCmd.Flags().IntVarP(&replicas, "replicas", "r", 0, "the number of replicas for the database")
	dbUpdateCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the database")
}
