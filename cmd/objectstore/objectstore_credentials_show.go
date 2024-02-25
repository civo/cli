package objectstore

import (
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var objectStoreCredentialFindCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "get"},
	Example: `civo objectstore credential show `,
	Short:   "Displays an Object Store Credential",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		creds, err := client.GetObjectStoreCredential(args[0])
		if err != nil {
			utility.Error("ObjectStore Credentials %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": creds.ID, "name": creds.Name, "Access Key ID": creds.AccessKeyID, "Secret AccessKey ID": creds.SecretAccessKeyID, "Max Size GB": strconv.Itoa(creds.MaxSizeGB), "Suspended": strconv.FormatBool(creds.Suspended), "Status": creds.Status})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
