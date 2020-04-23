package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var apikeyRemoveCmd = &cobra.Command{
	Use:     "remove NAME",
	Aliases: []string{"delete", "rm"},
	Short:   "Remove a saved API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		index, name, err := apiKeyFind(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		numKeys := len(config.Current)
		config.Current = append(config.Current[:index], config.Current[index+1:]...)
		config.SaveConfig()

		if numKeys > len(config.Current) {
			fmt.Printf("Removed the API Key %s\n", aurora.Green(name))
		} else {
			fmt.Fprintf(os.Stderr, "The API Key %s couldn't be found\n", aurora.Red(args[0]))
			os.Exit(1)
		}
	},
}
