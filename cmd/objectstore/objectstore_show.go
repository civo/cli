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

var objectStoreShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "info"},
	Example: `civo objectstore show OBJECTSTORE_NAME`,
	Short:   "Prints information about an Object Store",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		objectStore, err := client.FindObjectStore(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		ow.StartLine()
		fmt.Println()
		ow.AppendDataWithLabel("id", objectStore.ID, "ID")
		ow.AppendDataWithLabel("name", objectStore.Name, "Name")
		ow.AppendDataWithLabel("size", strconv.Itoa(objectStore.MaxSize), "Size")
		ow.AppendDataWithLabel("objectstore_endpoint", fmt.Sprintf("objectstore.%s.civo.com", strings.ToLower(client.Region)), "Object Store Endpoint")
		ow.AppendDataWithLabel("region", client.Region, "Region")
		ow.AppendDataWithLabel("accesskey", objectStore.OwnerInfo.AccessKeyID, "Access Key")
		ow.AppendDataWithLabel("status", objectStore.Status, "Status")

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
			fmt.Printf("To access the secret key run: civo objectstore credential secret --access-key=%s\n", utility.Green(objectStore.OwnerInfo.AccessKeyID))
		}
	},
}
