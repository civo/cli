package database

import (
	"errors"

	"github.com/spf13/cobra"
)

var name, schedule, backupType string
var scheduled bool

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
	dbBackupCmd.AddCommand(dbBackupShowCmd)
	dbBackupCmd.AddCommand(dbBackupDeleteCmd)

	// Create cmd options
	dbBackupCreateCmd.Flags().StringVarP(&name, "name", "n", "", "name of the database backup")
	dbBackupCreateCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "schedule of the database backup in the form of cronjob")
	dbBackupCreateCmd.Flags().StringVarP(&backupType, "type", "t", "scheduled", "set the type of database backup manual/scheduled")
	_ = dbBackupCreateCmd.MarkFlagRequired("name")
	addBackupDestinationFlags(dbBackupCreateCmd)

	// Update cmd options
	dbBackupUpdateCmd.Flags().StringVarP(&name, "name", "n", "", "name of the database backup")
	dbBackupUpdateCmd.Flags().StringVarP(&schedule, "schedule", "s", "", "schedule of the database backup in the form of cronjob")
	addBackupDestinationFlags(dbBackupUpdateCmd)

	// Delete cmd options
	dbBackupDeleteCmd.Flags().BoolVar(&scheduled, "scheduled", false, "delete scheduled backups instead of manual backups")
}

// addBackupDestinationFlags registers the flags that send a backup to the
// Civo Object Store of another region or to an external S3 bucket. Without
// them the backup stays in this region's Civo Object Store.
func addBackupDestinationFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&destinationRegion, "destination-region", "", "store the backup in the Civo Object Store of this region instead of the database's region")
	cmd.Flags().StringVar(&s3Endpoint, "s3-endpoint", "", "store the backup in an external S3-compatible bucket at this endpoint URL")
	cmd.Flags().StringVar(&s3Region, "s3-region", "", "region of the external S3 bucket")
	cmd.Flags().StringVar(&s3Bucket, "s3-bucket", "", "name of the external S3 bucket")
	cmd.Flags().StringVar(&s3AccessKey, "s3-access-key", "", "access key ID for the external S3 bucket")
	cmd.Flags().StringVar(&s3SecretKey, "s3-secret-key", "", "secret access key for the external S3 bucket (on update, leave empty to keep the current one)")
	cmd.Flags().StringVar(&s3PathPrefix, "s3-path-prefix", "", "optional path prefix inside the external S3 bucket")
}
