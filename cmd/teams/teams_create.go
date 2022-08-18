package teams

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var teamsCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo teams create NAME",
	Short:   "Create a new team",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		_, err = client.CreateTeam(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		team, err := client.FindTeam(args[0])
		if err != nil {
			utility.Error("Team %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": team.ID, "name": team.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Created a team called %s with team ID %s\n", utility.Green(team.Name), utility.Green(team.ID))
		}
	},
}
