package cmd

import (
	"fmt"
	"strings"

	"github.com/civo/cli/config"
	"github.com/spf13/cobra"
)

var apikeyCmd = &cobra.Command{
	Use:     "apikey",
	Aliases: []string{"apikeys"},
	Short:   "Manage API keys used to access your Civo account",
	Long: `If you use multiple Civo accounts, e.g. one for personal and one
for work, then you can setup multiple API keys and switch
between them when required.`,
}

func apiKeyFind(search string) (int, string, error) {
	var result int
	var name string
	for index, v := range config.Current {
		for key, value := range v.APIKeys {
			if strings.Contains(key, search) || strings.Contains(value, search) {
				result = index
				name = key
			}
		}
	}

	if name == "" {
		return 0, "", fmt.Errorf("unable to find %s at all in the list", search)
	}

	return result, name, nil
}

func init() {
	rootCmd.AddCommand(apikeyCmd)
	apikeyCmd.AddCommand(apikeyListCmd)
	apikeyCmd.AddCommand(apikeySaveCmd)
	apikeyCmd.AddCommand(apikeyRemoveCmd)
	apikeyCmd.AddCommand(apikeyCurrentCmd)
}
