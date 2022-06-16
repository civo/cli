package cmd

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var accessKey, name string

var objectStoreCredentialSecretCmd = &cobra.Command{
	Use:     "secret",
	Short:   "Access the secret key for the objectstore by providing your access key or by the objectstore name.",
	Example: "civo objectstore credential secret --access-key ACCESS_KEY / --name OBJECTSTORE_NAME",
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
		if name != "" {
			key = name
		}
		if key == "" {
			utility.Error("You must provide an access key or name")
			os.Exit(1)
		}

		objectStore, err := client.FindObjectStore(key)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if objectStore.Status == "creating" || objectStore.Status == "" {
			utility.Error("The objectstore is still being created. Please try again in a moment.")
			os.Exit(1)
		} else if objectStore.Status == "failed" {
			utility.Error("The objectstore failed to create. Please contact Civo support.")
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("accessKey", objectStore.AccessKeyID, "accessKey")
		ow.AppendDataWithLabel("secretKey", objectStore.SecretAccessKey, "secretKey")

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
