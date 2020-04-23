package cmd

import (
	"fmt"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var apikeySaveCmd = &cobra.Command{
	Use:     "save NAME KEY",
	Aliases: []string{"add", "store", "create", "save"},
	Short:   "Save a new API keys",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		apiKeyConfig := &config.Config{APIKeys: map[string]string{args[0]: args[1]}}
		config.Current = append(config.Current, *apiKeyConfig)
		config.SaveConfig()

		ow := utility.NewOutputWriterWithMap(map[string]string{"Name": args[0], "Key": args[1]})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Saved the API Key %s as %s\n", aurora.Green(args[0]), aurora.Green(args[1]))
		}
	},
}
