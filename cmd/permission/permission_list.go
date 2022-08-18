package permission

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var permissionsListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"all", "list"},
	Example: `civo permissions ls`,
	Short:   "List all available permissions",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
		}

		permissions, err := client.ListPermissions()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, permission := range permissions {
			ow.StartLine()

			ow.AppendDataWithLabel("code", permission.Code, "Code")
			ow.AppendDataWithLabel("name", permission.Name, "Name")
			ow.AppendDataWithLabel("description", permission.Description, "Description")
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
