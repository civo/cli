package objectstore

import (
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var format string

var objectStoreCredentialExportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"export-credentials"},
	Short:   "Export the credentials for your Object Store.",
	Example: "civo objectstore credential export --access-key=ACCESS_KEY --format=FORMAT (We support env and s3cfg)",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		var key string
		if accessKey != "" {
			if format == "" {
				utility.Error("You must provide a format to export to. See --help for more information.")
				os.Exit(1)
			}
			key = accessKey
		}
		if format != "" {
			if accessKey == "" {
				utility.Error("You must provide an access key. See --help for more information.")
				os.Exit(1)
			}
		}
		if key == "" {
			utility.Error("You must provide an access key and the format to export to. See --help for more information.")
			os.Exit(1)
		}

		credential, err := client.FindObjectStoreCredential(key)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if format == "env" {
			utility.Printf("# Tip: You can redirect output with (>> ~/.zshrc) to add these to Zsh's startup automatically\n")
			utility.Printf("export AWS_ACCESS_KEY_ID=%s\n", credential.AccessKeyID)
			utility.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", credential.SecretAccessKeyID)
			utility.Printf("export AWS_DEFAULT_REGION=%s\n", client.Region)
			utility.Printf("export AWS_HOST=https://objectstore.%s.civo.com\n", strings.ToLower(client.Region))
		} else if format == "s3cfg" {
			utility.Printf("# Tip: You can redirect output with (>> ~/.s3cfg) to automatically configure s3cmd\n")
			utility.Printf("[default]\n")
			utility.Printf("access_key = %s\n", credential.AccessKeyID)
			utility.Printf("secret_key = %s\n", credential.SecretAccessKeyID)
			utility.Printf("bucket_location = %s\n", client.Region)
			utility.Printf("host_base = objectstore.%s.civo.com\n", strings.ToLower(client.Region))
			utility.Printf("signature_v2 = True")
		} else {
			utility.Error("You must provide a valid format to export to. Supported formats are env and s3cfg. See --help for more information.")
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		}
	},
}
