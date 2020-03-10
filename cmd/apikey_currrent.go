package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var apikeyCurrentCmd = &cobra.Command{
	Use:     "current [NAME]",
	Aliases: []string{"use", "default", "set"},
	Short:   "Show the current API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, err := apiKeyFind(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		value := config.Current.APIKeys[name]
		if value != "" {
			ow := utility.NewOutputWriterWithMap(map[string]string{"Name": name, "Key": value})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("Set the default API Key to be %s\n", aurora.Green(name))
			}
		}
	},
}
