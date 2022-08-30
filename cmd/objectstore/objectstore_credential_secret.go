package objectstore

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var accessKey string

var objectStoreCredentialSecretCmd = &cobra.Command{
	Use:     "secret",
	Short:   "Access the secret key for the Object Store by providing your access key.",
	Example: "civo objectstore credential secret --access-key ACCESS_KEY",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		var key string
		if accessKey != "" {
			key = accessKey
		}
		if key == "" {
			utility.Error("You must provide an access key. See --help for more information.")
			os.Exit(1)
		}

		credential, err := client.FindObjectStoreCredential(key)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if credential.Status == "creating" || credential.Status == "" {
			utility.Error("The Object Store Credential is still being created. Please try again in a moment.")
			os.Exit(1)
		} else if credential.Status == "failed" {
			utility.Error("The Object Store Credential failed to create. Please contact Civo support.")
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Your secret key is: %s\n", utility.Green(credential.SecretAccessKeyID))
		}
	},
}
