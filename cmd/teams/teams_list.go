package teams

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var teamsListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"all", "list"},
	Example: `civo teams ls`,
	Short:   "List all teams",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
		}
		teams, err := client.ListTeams()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, team := range teams {
			ow.StartLine()

			ow.AppendDataWithLabel("id", team.ID, "ID")
			ow.AppendDataWithLabel("name", team.Name, "Name")
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
