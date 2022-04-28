package cmd

import (
	"errors"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var appDomainList []string
var appDomainRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove a domain for your application.",
	Example: "civo app domain rm APP_NAME DOMAIN_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findApp, err := client.FindApplication(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s application in your account", utility.Red(args[0]))
				os.Exit(1)
			} else if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one application with that name in your account")
				os.Exit(1)
			} else {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		for _, appDomain := range findApp.Domains {
			if appDomain == args[1] {
				remove(findApp.Domains, appDomain)
				utility.Info("Domain %s removed from %s", utility.Green(appDomain), utility.Green(args[0]))
				break
			}
		}
	},
}

func remove(items []string, item string) []string {
	newitems := []string{}

	for _, i := range items {
		if i != item {
			newitems = append(newitems, i)
		}
	}

	return newitems
}
