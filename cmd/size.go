package cmd

import (
	"github.com/spf13/cobra"
)

var sizeCmd = &cobra.Command{
	Use:     "size",
	Aliases: []string{"sizes"},
	Short:   "Details of Civo instance sizes",
}

func init() {
	rootCmd.AddCommand(sizeCmd)
	sizeCmd.AddCommand(sizeListCmd)
}
