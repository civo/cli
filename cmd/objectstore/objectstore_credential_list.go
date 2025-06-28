package objectstore

import (
	"os"

	"github.com/civo/civogo"
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
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		var creds []civogo.ObjectStoreCredential
		page := 1
		perPage := 100

		for {
			paginatedCreds, err := client.ListObjectStoreCredentials(page, perPage)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			creds = append(creds, paginatedCreds.Items...)
			if page >= paginatedCreds.Pages {
				break
			}
			page++
		}

		ow := utility.NewOutputWriter()
		for _, credential := range creds {
			ow.StartLine()
			ow.AppendDataWithLabel("name", credential.Name, "Name")
			ow.AppendDataWithLabel("access_key", credential.AccessKeyID, "Access Key")
			ow.AppendDataWithLabel("status", credential.Status, "Status")
		}

		ow.FinishAndPrintOutput()
	},
}
