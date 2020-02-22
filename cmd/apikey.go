package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// apikeyCmd represents the apikey command
var apikeyCmd = &cobra.Command{
	Use:     "apikey",
	Aliases: []string{"apikeys"},
	Short:   "Manage API keys used to access your Civo account",
	Long: `If you use multiple Civo accounts, e.g. one for personal and one
for work, then you can setup multiple API keys and switch
between them when required.`,
}

// apikeyListCmd represents the command to list available API keys
var apikeyListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List stored API keys",
	Long: `List all API keys, making clear which is the current default.
If you wish to use a custom format, the available fields are:

* Name
* Key

Example: civo apikey ls -o custom -f "Name: Key"`,
	Run: func(cmd *cobra.Command, args []string) {
		data := make([][]string, len(CurrentConfig.APIKeys))

		key := 0
		for name, apiKey := range CurrentConfig.APIKeys {
			defaultLabel := ""
			if CurrentConfig.Meta.CurrentAPIKey == name {
				defaultLabel = "<====="
			}
			data[key] = []string{name, apiKey, defaultLabel}
			key++
		}

		outputTable([]string{"Name", "Key", "Default"}, data)
	},
}

// apikeySaveCmd represents the command to save a new API key
var apikeySaveCmd = &cobra.Command{
	Use:     "save NAME KEY",
	Aliases: []string{"add", "store", "create"},
	Short:   "Save a new API keys",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		CurrentConfig.APIKeys[args[0]] = args[1]
		saveConfig()

		if OutputFormat == "human" {
			fmt.Printf("Saved the API Key %s as %s\n", aurora.Green(args[0]), aurora.Green(args[1]))
		} else {
			outputKeyValue(map[string]string{"Name": args[0], "Key": args[1]})
		}
	},
}

// apikeyRemoveCmd represents the command to remove a saved API key
var apikeyRemoveCmd = &cobra.Command{
	Use:     "remove NAME",
	Aliases: []string{"delete", "rm"},
	Short:   "Remove a saved API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		numKeys := len(CurrentConfig.APIKeys)
		delete(CurrentConfig.APIKeys, args[0])
		saveConfig()

		if numKeys > len(CurrentConfig.APIKeys) {
			fmt.Printf("Removed the API Key %s\n", aurora.Green(args[0]))
		} else {
			fmt.Printf("The API Key %s couldn't be found\n", aurora.Red(args[0]))
			os.Exit(1)
		}
	},
}

// apikeyCurrentCmd represents the command to show the current API key
var apikeyCurrentCmd = &cobra.Command{
	Use:     "current [NAME]",
	Aliases: []string{"use", "default"},
	Short:   "Show the current API key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// If specified
		//   Find based on single partial match
		//   Set
		//   Output
		// else
		//   Output
		value := CurrentConfig.APIKeys[args[0]]
		if value != "" {
			if OutputFormat == "human" {
				fmt.Printf("Saved the API Key %s as %s\n", aurora.Green(args[0]), aurora.Green(args[1]))
			} else {
				outputKeyValue(map[string]string{"Name": args[0], "Key": args[1]})
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(apikeyCmd)

	apikeyCmd.AddCommand(apikeyListCmd)

	apikeyCmd.AddCommand(apikeySaveCmd)

	apikeyCmd.AddCommand(apikeyRemoveCmd)

	apikeyCmd.AddCommand(apikeyCurrentCmd)
}
