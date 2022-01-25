package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var teamsRenameCmd = &cobra.Command{
	Use:     "rename",
	Short:   "Rename a team",
	Example: "civo teams rename OLD_TEAM_NAME NEW_TEAM_NAME",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()

		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findTeam, err := client.FindTeam(args[0])
		if err != nil {
			utility.Error("Team %s", err)
			os.Exit(1)
		}

		newTeamName := args[1]

		team, err := client.RenameTeam(findTeam.ID, newTeamName)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": team.ID, "name": team.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The team with ID %s was renamed to %s\n", utility.Green(team.ID), utility.Green(team.Name))
		}
	},
}
