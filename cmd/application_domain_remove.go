package cmd

import (
	"errors"
	"fmt"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var domainsToBeDeleted []string
var appDomainRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove a domain for your application.",
	Example: "civo app domain rm APP_NAME DOMAIN_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
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

		if utility.UserConfirmedDeletion(fmt.Sprintf("domain %s", pluralize.Pluralize(len(appList), "")), defaultYes, strings.Join(domainsToBeDeleted, ", ")) {

			for _, appDomain := range findApp.Domains {
				if appDomain == args[1] {
					domainsToBeDeleted = remove(findApp.Domains, appDomain)
				}
			}

			application := &civogo.UpdateApplicationRequest{
				Domains: domainsToBeDeleted,
			}

			app, err := client.UpdateApplication(findApp.ID, application)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			ow := utility.NewOutputWriterWithMap(map[string]string{"id": app.ID, "name": app.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON(prettySet)
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("\nThe domain %s has been removed from your application %s \n", utility.Green(args[1]), utility.Green(app.Name))

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
