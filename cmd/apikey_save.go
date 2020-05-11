package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var apikeySaveCmd = &cobra.Command{
	Use:     "save NAME KEY",
	Aliases: []string{"add", "store", "create", "save"},
	Short:   "Save a new API keys",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		config.Current.APIKeys[args[0]] = args[1]
		config.SaveConfig()

		ow := utility.NewOutputWriterWithMap(map[string]string{"Name": args[0], "Key": args[1]})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Saved the API Key %s as %s\n", utility.Green(args[0]), utility.Green(args[1]))
		}
	},
}
