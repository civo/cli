package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
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
		keys := make([]string, 0, len(config.Current.APIKeys))
		for k := range config.Current.APIKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		ow := utility.NewOutputWriter()

		for _, name := range keys {
			ow.StartLine()
			apiKey := config.Current.APIKeys[name]
			defaultLabel := ""
			if config.Current.Meta.CurrentAPIKey == name {
				defaultLabel = "<====="
			}
			ow.AppendData("Name", name)
			ow.AppendData("Key", apiKey)
			ow.AppendData("Default", defaultLabel)
		}

		switch OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(OutputFields)
		default:
			ow.WriteTable()
		}
	},
}

// apikeySaveCmd represents the command to save a new API key
var apikeySaveCmd = &cobra.Command{
	Use:     "save NAME KEY",
	Aliases: []string{"add", "store", "create", "save"},
	Short:   "Save a new API keys",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		config.Current.APIKeys[args[0]] = args[1]
		config.SaveConfig()

		ow := utility.NewOutputWriterWithMap(map[string]string{"Name": args[0], "Key": args[1]})

		switch OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(OutputFields)
		default:
			fmt.Printf("Saved the API Key %s as %s\n", aurora.Green(args[0]), aurora.Green(args[1]))
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
		key, err := apiKeyFind(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		numKeys := len(config.Current.APIKeys)
		delete(config.Current.APIKeys, key)
		config.SaveConfig()

		if numKeys > len(config.Current.APIKeys) {
			fmt.Printf("Removed the API Key %s\n", aurora.Green(args[0]))
		} else {
			fmt.Fprintf(os.Stderr, "The API Key %s couldn't be found\n", aurora.Red(args[0]))
			os.Exit(1)
		}
	},
}

// apikeyCurrentCmd represents the command to show the current API key
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

			switch OutputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(OutputFields)
			default:
				fmt.Printf("Set the default API Key to be %s\n", aurora.Green(name))
			}
		}
	},
}

func apiKeyFind(search string) (string, error) {
	var result string
	for k, v := range config.Current.APIKeys {
		if strings.Contains(k, search) || strings.Contains(v, search) {
			if result == "" {
				result = k
			} else {
				return "", fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
		}
	}

	if result == "" {
		return "", fmt.Errorf("unable to find %s at all in the list", search)
	}

	return result, nil
}

func init() {
	rootCmd.AddCommand(apikeyCmd)

	apikeyCmd.AddCommand(apikeyListCmd)

	apikeyCmd.AddCommand(apikeySaveCmd)

	apikeyCmd.AddCommand(apikeyRemoveCmd)

	apikeyCmd.AddCommand(apikeyCurrentCmd)
}
