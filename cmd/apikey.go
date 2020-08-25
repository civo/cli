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

func apiKeyFind(search string) (string, error) {
	var result string
	for key, value := range config.Current.APIKeys {
		if strings.Contains(key, search) || strings.Contains(value, search) {
			result = key
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
	apikeyCmd.AddCommand(apikeyShowCmd)
}
