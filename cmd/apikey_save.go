package cmd

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var apikeySaveCmdExample = `* Interactive way:
    civo apikey save

* Non-interactive way:
    civo apikey save NAME APIKEY

* Load from environment variables way:
    civo apikey save --load-from-env

Notes:
* This command will generate one file called '.civo.json' in your home directory
* The NAME is just an identifier for your own reference. This can be useful if you have multiple accounts.
* Some ideas for NAME are: 'personal', 'work', 'ci-server', 'staging', 'production'
* When --load-from-env flag is provided, we assume you have set the following environment variables:
  * (required) CIVO_API_KEY e.g. 'export CIVO_API_KEY=<YOUR_CIVO_API_KEY>'
  * (optional) CIVO_API_KEY_NAME e.g. 'export CIVO_API_KEY_NAME=personal'
  * When CIVO_API_KEY_NAME is not set, it will default to the hostname where the this CLI is running
`

var loadAPIKeyFromEnv bool

var apikeySaveCmd = &cobra.Command{
	Use:     "save",
	Aliases: []string{"add", "store", "create", "new"},
	Short:   "Save a new API key",
	Example: apikeySaveCmdExample,
	Run: func(cmd *cobra.Command, args []string) {

		var name, apiKey string
		var err error

		// if arg is more than two, return an error
		if len(args) > 2 {
			utility.Info("There are too many arguments for this command")
			cmd.Help()
			os.Exit(1)
		}

		if len(args) == 0 && !loadAPIKeyFromEnv {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter a nice name for this account/API Key: ")

			name, err = reader.ReadString('\n')
			if err != nil {
				utility.Error("Error reading name", err)
				os.Exit(1)
			}
			if runtime.GOOS == "windows" {
				name = strings.TrimSuffix(name, "\r\n")
			} else {
				name = strings.TrimSuffix(name, "\n")
			}
			fmt.Printf("Enter the API key: ")
			apikeyBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				utility.Error("Error reading api key", err)
				os.Exit(1)
			}
			apiKey = string(apikeyBytes)
		}

		if len(args) == 2 && !loadAPIKeyFromEnv {
			name = args[0]
			apiKey = args[1]
		}

		if loadAPIKeyFromEnv {
			nameEnvRef := "CIVO_API_KEY_NAME"
			nameEnv, present := os.LookupEnv(nameEnvRef)
			if !present || nameEnv == "" {
				hostname, err := os.Hostname()
				if err != nil {
					utility.Error("unable to retrieve hostname - %s", err)
					os.Exit(1)
				}
				nameEnv = hostname
			}

			apiKeyEnvRef := "CIVO_API_KEY"
			apiKeyEnv, present := os.LookupEnv(apiKeyEnvRef)
			if !present || apiKeyEnv == "" {
				utility.Error("%q environment variable is missing", apiKeyEnvRef)
				os.Exit(1)
			}

			name = nameEnv
			apiKey = apiKeyEnv
		}

		config.Current.APIKeys[name] = apiKey
		if config.Current.Meta.DefaultRegion == "" {
			client, err := civogo.NewClientWithURL(apiKey, config.Current.Meta.URL, "")
			if err != nil {
				utility.Error("Unable to create a Civo API client, please report this at https://github.com/civo/cli")
				os.Exit(1)
			}
			region, err := client.GetDefaultRegion()
			if err != nil {
				utility.Error("Unable to get the default regions from the Civo API")
				os.Exit(1)
			}
			config.Current.Meta.DefaultRegion = region.Code
		}

		config.SaveConfig()

		if len(config.Current.APIKeys) == 1 {
			config.Current.Meta.CurrentAPIKey = name
			config.SaveConfig()
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"name": name, "key": apiKey})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Saved the API Key %s\n", utility.Green(name))
		}

	},
}
