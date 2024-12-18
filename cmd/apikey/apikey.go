package apikey

import (
	"errors"
	"fmt"

	"github.com/civo/cli/config"
	"github.com/spf13/cobra"
)

// APIKeyCmd is the root command for `civo apikey`
var APIKeyCmd = &cobra.Command{
	Use:     "apikey",
	Aliases: []string{"apikeys"},
	Short:   "Manage API keys used to access your Civo account",
	Long: `If you use multiple Civo accounts, e.g. one for personal and one
for work, then you can setup multiple API keys and switch
between them when required.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func apiKeyFind(search string) (string, error) {
	var result string
	for key, value := range config.Current.APIKeys {
		if key == search || value == search {
			result = key
		}
	}

	if result == "" {
		return "", fmt.Errorf("unable to find %s at all in the list", search)
	}

	return result, nil
}

func init() {

	config.SkipAPIInitialization = true

	APIKeyCmd.AddCommand(apikeyListCmd)
	APIKeyCmd.AddCommand(apikeySaveCmd)
	APIKeyCmd.AddCommand(apikeyRemoveCmd)
	APIKeyCmd.AddCommand(apikeyCurrentCmd)
	APIKeyCmd.AddCommand(apikeyShowCmd)

	// Flags for "civo apikey save" command
	apikeySaveCmd.Flags().BoolVar(&loadAPIKeyFromEnv, "load-from-env", false, "When set, the name and key will be taken from environment variables (see notes above)")
	apikeyRemoveCmd.Flags().BoolVar(&forceFlag, "force", false, "Force removal of the current API key without confirmation")
}
