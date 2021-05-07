package cmd

import (
	"fmt"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
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

func setCommandFlags(c *cobra.Command) {
	c.Flags().StringVarP(&config.Filename, "config", "", "", "config file (default is $HOME/.civo.json)")
	c.Flags().StringVarP(&outputFields, "fields", "f", "", "output fields for custom format output (use -h to determine fields)")
	c.Flags().StringVarP(&outputFormat, "output", "o", "human", "output format (json/human/custom)")
	c.Flags().BoolVarP(&defaultYes, "yes", "y", false, "Automatic yes to prompts; assume \"yes\" as answer to all prompts and run non-interactively")
	c.Flags().StringVarP(&regionSet, "region", "", "", "Choose the region to connect to, if you use this option it will use it over the default region")
}

func searchAllCommands(allCommands *[]*cobra.Command, commands []*cobra.Command) {
	for _, c := range commands {
		*allCommands = append(*allCommands, c)
		searchAllCommands(allCommands, c.Commands())
	}
}

func isKubemart(commandPath string) bool {
	splitted := strings.Fields(commandPath)
	return splitted[1] == kubemartCmd.Name()
}

func init() {
	cobra.OnInitialize(config.ReadConfig)
	setCommandFlags(rootCmd)

	// Set Civo CLI flags to all commands except "kubemart" and its sub-commands.
	// Show only Kubemart flags when running "civo kubemart" and "civo kubemart <command>".
	allCommands := []*cobra.Command{}
	searchAllCommands(&allCommands, rootCmd.Commands())
	for _, c := range allCommands {
		if !isKubemart(c.CommandPath()) {
			setCommandFlags(c)
		}
	}

	// Add warning if the region is empty, for the user with the old config
	config.ReadConfig()
	if config.Current.Meta.DefaultRegion == "" {
		utility.Warning("No region set - using the default one - set a default using \"civo region current REGION\" or specify one with every command using \"--region=REGION\"")
	}

}
