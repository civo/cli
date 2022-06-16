package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var objectStoreCmd = &cobra.Command{
	Use:     "objectstore",
	Aliases: []string{"bucket"},
	Short:   "Civo Objectstore/Bucket management commands",
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
	Aliases: []string{"credentials"},
	Short:   "Credentials for Civo objectstore",
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
	objectStoreCreateCmd.Flags().IntVarP(&maxObjects, "max-objects", "m", 1000, "Maximum number of objects in the bucket")

	//Flags for update cmd
	objectStoreUpdateCmd.Flags().IntVarP(&bucketSize, "size", "s", 500, "Size of the bucket")
	objectStoreUpdateCmd.Flags().IntVarP(&maxObjects, "max-objects", "m", 1000, "Maximum number of objects in the bucket")

	//Credential commands
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialSecretCmd)
	objectStoreCredentialSecretCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialSecretCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the objectstore")
}
