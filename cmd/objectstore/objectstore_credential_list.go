package objectstore

import (
	"os"
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var objectStoreCredentialListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List all Object Store Credentials",
	Example: "civo objectstore credential ls",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		creds, err := client.ListObjectStoreCredentials()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, credential := range creds.Items {
			ow.StartLine()
			ow.AppendDataWithLabel("id", credential.ID, "ID")
			ow.AppendDataWithLabel("name", credential.Name, "Name")
			ow.AppendDataWithLabel("size", strconv.Itoa(credential.MaxSizeGB), "Size")
			ow.AppendDataWithLabel("status", credential.Status, "Status")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
