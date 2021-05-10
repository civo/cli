package cmd

import (
	"fmt"

	"github.com/civo/cli/config"
	"github.com/spf13/cobra"
)

var outputFields, outputFormat, regionSet string
var defaultYes bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "civo",
	Short: "CLI to manage cloud resources at Civo.com",
	Long: `civo is a CLI library for managing cloud resources such
as instances and Kubernetes clusters at Civo.com.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		// if len(config.Current.APIKeys) == 0 {
		// 	utility.Warning("You need to add a api key")
		// 	utility.Info("1 - Open https://www.civo.com/account/security")
		// 	utility.Info("2 - Copy the API key")
		// 	utility.Info("3 - Run civo apikey add NAME API_KEY")
		// 	os.Exit(1)
		// }

		// if config.DefaultAPIKey() == "" {
		// 	utility.Warning("You need to define a default api key")
		// 	os.Exit(1)
		// }
		fmt.Println(err)
		// os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.ReadConfig)

	rootCmd.PersistentFlags().StringVarP(&config.Filename, "config", "", "", "config file (default is $HOME/.civo.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFields, "fields", "f", "", "output fields for custom format output (use -h to determine fields)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "human", "output format (json/human/custom)")
	rootCmd.PersistentFlags().BoolVarP(&defaultYes, "yes", "y", false, "Automatic yes to prompts; assume \"yes\" as answer to all prompts and run non-interactively")
	rootCmd.PersistentFlags().StringVarP(&regionSet, "region", "", "", "Choose the region to connect to, if you use this option it will use it over the default region")

	// Add warning if the region is empty, for the user with the old config
	config.ReadConfig()
}
