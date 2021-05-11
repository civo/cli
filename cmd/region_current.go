package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var regionCurrentCmd = &cobra.Command{
	Use:     "current [NAME]",
	Aliases: []string{"use", "default", "set"},
	Short:   "Set the current region",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			return
		}

		regions, err := client.ListRegions()
		if err != nil {
			utility.Error("%s", err)
			return
		}

		if config.Current.Meta.DefaultRegion == args[0] {
			fmt.Printf("You are already using that region: %s\n", utility.Red(args[0]))
			os.Exit(1)
		}

		for _, v := range regions {
			if v.Code == args[0] {
				config.Current.Meta.DefaultRegion = args[0]
				config.SaveConfig()
				fmt.Printf("The default region was set to (%s) %s\n", v.Name, utility.Green(args[0]))
				os.Exit(1)
			}
		}

		fmt.Printf("The region you tried to set %s doesn't exist, please use 'civo region ls' to get the code of a valid region", utility.Red(args[0]))
		os.Exit(1)
	},
}
