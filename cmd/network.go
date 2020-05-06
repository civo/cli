package cmd

// network list -- list all networks [ls, all]
// network create LABEL -- create a new private network called LABEL [new]
// network update OLD_NAME NEW_NAME -- create a new private network called LABEL [new]
// network remove ID -- remove the network ID [delete, destroy, rm]

import (
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"net", "nw"},
	Short:   "Details of Civo Network",
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(networkListCmd)
	networkCmd.AddCommand(networkCreateCmd)
	networkCmd.AddCommand(networkUpdateCmd)
	networkCmd.AddCommand(networkRemoveCmd)
}
