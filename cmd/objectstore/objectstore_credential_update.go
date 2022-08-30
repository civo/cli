package objectstore

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var credAccessKey, credSecretAccessKey string
var credSuspended bool

var objectStoreCredentialUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit", "modify", "change", "update"},
	Short:   "Update an Object Store Credential",
	Example: "civo objectstore credential update CREDENTIAL_NAME --size SIZE",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		credential, err := client.FindObjectStoreCredential(args[0])
		if err != nil {
			utility.Error("Object Store Credential %s", err)
			os.Exit(1)
		}

		if credentialSize != 0 {
			credential.MaxSizeGB = credentialSize
		}
		if credAccessKey != "" {
			credential.AccessKeyID = credAccessKey
		}
		if credSecretAccessKey != "" {
			credential.SecretAccessKeyID = credSecretAccessKey
		}
		if credSuspended {
			credential.Suspended = true
		}

		cred, err := client.UpdateObjectStoreCredential(credential.ID, &civogo.UpdateObjectStoreCredentialRequest{
			MaxSizeGB:         &credential.MaxSizeGB,
			AccessKeyID:       &credential.AccessKeyID,
			SecretAccessKeyID: &credential.SecretAccessKeyID,
			Suspended:         credential.Suspended,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": cred.ID, "name": cred.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The Object Store Credential with ID %s is updated\n", utility.Green(cred.ID))
			os.Exit(0)
		}
	},
}
