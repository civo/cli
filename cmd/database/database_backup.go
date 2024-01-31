package database

import (
	"errors"

	"github.com/spf13/cobra"
)

var name, schedule, backupType string

// dbBackupCmd is the root command for the db backup subcommand
var dbBackupCmd = &cobra.Command{
	Use:     "backup",
	Aliases: []string{"bk", "backups"},
	Short:   "Manage Civo Database Backups",
	Long:    `Create, update, and list Civo Database Backups.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	dbBackupCmd.AddCommand(dbBackupCreateCmd)
	dbBackupCmd.AddCommand(dbBackupListCmd)
	dbBackupCmd.AddCommand(dbBackupUpdateCmd)

	// Create cmd options
	dbBackupCreateCmd.Flags().StringVarP(&name, "name", "n", "", "name of the database backup")
	dbBackupCreateCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "schedule of the database backup in the form of cronjob")
	dbBackupCreateCmd.Flags().StringVarP(&backupType, "type", "t", "scheduled", "set the type of database backup manual/scheduled")

	dbBackupCreateCmd.MarkFlagRequired("name")

	// Update cmd options
	dbBackupUpdateCmd.Flags().StringVarP(&name, "name", "n", "", "name of the database backup")
	dbBackupUpdateCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "schedule of the database backup in the form of cronjob")
}
