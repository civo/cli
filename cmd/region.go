package cmd

import (
	"github.com/spf13/cobra"
)

var regionCmd = &cobra.Command{
	Use:     "region",
	Aliases: []string{"regions"},
	Short:   "Details of Civo regions",
}

func init() {
	rootCmd.AddCommand(regionCmd)
	regionCmd.AddCommand(regionListCmd)
	regionCmd.AddCommand(regionCurrentCmd)
}
