package objectstore

import (
	"errors"

	"github.com/spf13/cobra"
)

var credentialSize int
var accessKey, secretAccessKey string

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
	objectStoreCreateCmd.Flags().Int64VarP(&bucketSize, "size", "s", 500, "Size of the Object store")
	objectStoreCreateCmd.Flags().StringVarP(&owner, "owner-name", "n", "", "Name of Owner of the Object store. You can reference name of any civo object store credential created before")
	objectStoreCreateCmd.Flags().StringVarP(&owner, "owner-access-key", "a", "", "Access Key ID of Owner of the Object store. You can reference name of any civo object store credential created before")
	objectStoreCreateCmd.Flags().BoolVarP(&waitOS, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the Object Store to be ready")

	//Flags for update cmd
	objectStoreUpdateCmd.Flags().Int64VarP(&bucketSize, "size", "s", 500, "Size of the object store")

	//Credential commands
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialSecretCmd)
	objectStoreCredentialSecretCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialSecretCmd.MarkFlagRequired("access-key")
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialExportCmd)
	objectStoreCredentialExportCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialExportCmd.Flags().StringVarP(&format, "format", "", "env", "Format of the output (We support env and s3cfg formats.)")
	objectStoreCredentialExportCmd.MarkFlagRequired("access-key")
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialListCmd)
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialCreateCmd)
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialUpdateCmd)
	objectStoreCredentialCmd.AddCommand(objectStoreCredentialDeleteCmd)

	//Flags for credential create command
	objectStoreCredentialCreateCmd.Flags().BoolVarP(&waitOS, "wait", "w", false, "a simple flag (e.g. --wait) that will cause the CLI to spin and wait for the credential to be ready")
	objectStoreCredentialCreateCmd.Flags().StringVarP(&accessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialCreateCmd.Flags().StringVarP(&secretAccessKey, "secret-access-key", "s", "", "Secret Access Key")

	//Flags for credential update command
	objectStoreCredentialUpdateCmd.Flags().StringVarP(&credAccessKey, "access-key", "a", "", "Access Key")
	objectStoreCredentialUpdateCmd.Flags().StringVarP(&credSecretAccessKey, "secret-key", "k", "", "Secret Key")
}
