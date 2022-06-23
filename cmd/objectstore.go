package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var objectStoreCmd = &cobra.Command{
	Use:     "objectstore",
	Aliases: []string{"bucket", "buckets", "object", "objects"},
	Short:   "Civo Object Store/Bucket management commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var objectStoreCredentialCmd = &cobra.Command{
	Use:     "credential",
	Aliases: []string{"credentials", "creds", "user", "users", "key", "keys"},
	Short:   "Credentials for Civo Object Store",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("command is required")
	},
}

func init() {
	rootCmd.AddCommand(objectStoreCmd)

	objectStoreCmd.AddCommand(objectStoreListCmd)
	objectStoreCmd.AddCommand(objectStoreCreateCmd)
	objectStoreCmd.AddCommand(objectStoreUpdateCmd)
	objectStoreCmd.AddCommand(objectStoreDeleteCmd)
	objectStoreCmd.AddCommand(objectStoreShowCmd)
	objectStoreCmd.AddCommand(objectStoreCredentialCmd)

	//Flags for create cmd
	objectStoreCreateCmd.Flags().IntVarP(&bucketSize, "size", "s", 500, "Size of the bucket")
	objectStoreCreateCmd.Flags().BoolVarP(&waitOS, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the Object Store to be ready")

	//Flags for update cmd
	objectStoreUpdateCmd.Flags().IntVarP(&bucketSize, "size", "s", 500, "Size of the bucket")

	//Credential commands
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialSecretCmd)
	objectStoreCredentialSecretCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialExportCmd)
	objectStoreCredentialExportCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialExportCmd.Flags().StringVarP(&format, "format", "", "env", "Format of the output (We support env and s3cfg formats.)")
}
