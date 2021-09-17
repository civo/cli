package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var diskImageCmd = &cobra.Command{
	Use:     "diskimage",
	Aliases: []string{"diskimages", "template", "templates"},
	Short:   "Details of Civo disk images",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	rootCmd.AddCommand(diskImageCmd)
	diskImageCmd.AddCommand(diskImageListCmd)
}
