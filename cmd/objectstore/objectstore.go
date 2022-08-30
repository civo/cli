package objectstore

import (
	"errors"

	"github.com/spf13/cobra"
)

//ObjectStoreCmd manages Civo Object Store
var ObjectStoreCmd = &cobra.Command{
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

	ObjectStoreCmd.AddCommand(objectStoreListCmd)
	ObjectStoreCmd.AddCommand(objectStoreCreateCmd)
	ObjectStoreCmd.AddCommand(objectStoreUpdateCmd)
	ObjectStoreCmd.AddCommand(objectStoreDeleteCmd)
	ObjectStoreCmd.AddCommand(objectStoreShowCmd)
	ObjectStoreCmd.AddCommand(objectStoreCredentialCmd)

	//Flags for create cmd
	objectStoreCreateCmd.Flags().Int64VarP(&bucketSize, "size", "s", 500, "Size of the bucket")
	objectStoreCreateCmd.Flags().StringVarP(&owner, "owner", "", "", "Owner of the bucket")
	objectStoreCreateCmd.Flags().BoolVarP(&waitOS, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the Object Store to be ready")

	//Flags for update cmd
	objectStoreUpdateCmd.Flags().Int64VarP(&bucketSize, "size", "s", 500, "Size of the bucket")

	//Credential commands
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialSecretCmd)
	objectStoreCredentialSecretCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialExportCmd)
	objectStoreCredentialExportCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialExportCmd.Flags().StringVarP(&format, "format", "", "env", "Format of the output (We support env and s3cfg formats.)")
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialListCmd)
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialCreateCmd)
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialUpdateCmd)
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialDeleteCmd)

	//Flags for credential create command
	objectStoreCredentialCreateCmd.Flags().IntVarP(&credentialSize, "size", "s", 500, "Size to allocate to the credential")
	objectStoreCredentialCreateCmd.Flags().BoolVarP(&waitOS, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the credential to be ready")

	//Flags for credential update command
	objectStoreCredentialUpdateCmd.Flags().IntVarP(&credentialSize, "size", "s", 500, "Size to update to")
	objectStoreCredentialUpdateCmd.Flags().StringVarP(&credAccessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialUpdateCmd.Flags().StringVarP(&credSecretAccessKey, "secret-key", "k", "", "Secret Key")
	objectStoreCredentialUpdateCmd.Flags().BoolVarP(&credSuspended, "suspended", "u", false, "Suspend the credential")
}
