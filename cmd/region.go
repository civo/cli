package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var regionCmd = &cobra.Command{
	Use:     "region",
	Aliases: []string{"regions"},
	Short:   "Details of Civo regions",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	rootCmd.AddCommand(regionCmd)
	regionCmd.AddCommand(regionListCmd)
	regionCmd.AddCommand(regionCurrentCmd)
}
