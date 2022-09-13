package objectstore

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var objectStoreListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo objectstore ls`,
	Short:   "List all Object Stores",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		objectStores, err := client.ListObjectStores()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, objectStore := range objectStores.Items {
			ow.StartLine()
			ow.AppendDataWithLabel("id", objectStore.ID[:6], "ID")
			ow.AppendDataWithLabel("name", objectStore.Name, "Name")
			ow.AppendDataWithLabel("max_size", strconv.Itoa(objectStore.MaxSize), "Size")
			ow.AppendDataWithLabel("objectstore_endpoint", fmt.Sprintf("objectstore.%s.civo.com", strings.ToLower(client.Region)), "Object Store Endpoint")
			ow.AppendDataWithLabel("status", objectStore.Status, "Status")
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
