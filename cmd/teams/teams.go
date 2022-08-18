package teams

import (
	"errors"

	"github.com/spf13/cobra"
)

// TeamsCmd manages Civo teams
var TeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Manage teams in Civo",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	TeamsCmd.AddCommand(teamsListCmd)
	TeamsCmd.AddCommand(teamsCreateCmd)
	TeamsCmd.AddCommand(teamsRenameCmd)
	TeamsCmd.AddCommand(teamsDeleteCmd)
}
