package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var teamList []utility.ObjecteList
var teamsDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"delete", "rm"},
	Short:   "Delete a team",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo teams delete TEAM_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}
		if len(args) == 1 {
			team, err := client.FindTeam(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s team in your organisation", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one team with that name in your organisation")
					os.Exit(1)
				}
			}
			teamList = append(teamList, utility.ObjecteList{ID: team.ID, Name: team.Name})
		} else {
			for _, v := range args {
				team, err := client.FindTeam(v)
				if err == nil {
					teamList = append(teamList, utility.ObjecteList{ID: team.ID, Name: team.Name})
				}
			}
		}

		teamNameList := []string{}
		for _, v := range teamList {
			teamNameList = append(teamNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("team %s", pluralize.Pluralize(len(teamList), "")), common.DefaultYes, strings.Join(teamNameList, ", ")) {

			for _, v := range teamList {
				_, err = client.DeleteTeam(v.ID)
				if err != nil {
					utility.Error("error deleting team: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range teamList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(teamList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The team %s(%s) has been deleted\n", pluralize.Pluralize(len(teamList), ""), utility.Green(strings.Join(teamNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
