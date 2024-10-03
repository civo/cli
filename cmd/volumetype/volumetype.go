package volumetype

import (
	"errors"
	"github.com/spf13/cobra"
)

var VolumeTypeCmd = &cobra.Command{
	Use:     "volumetypes",
	Aliases: []string{"voltype", "volumetype"},
	Short:   "Details of Civo Volume Types",
	Long:    `Commands to manage volume types in Civo cloud`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	VolumeTypeCmd.AddCommand(volumetypesListCmd)
}
