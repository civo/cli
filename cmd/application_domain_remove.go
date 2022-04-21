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
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one application with that name in your account")
				os.Exit(1)
			}
		}

		if len(args) == 2 {
			domain, err := client.FindAppDomain(args[1], findApp.ID)
			if err != nil {
				if errors.Is(err, civogo.ErrDNSDomainNotFound) {
					utility.Error("sorry there is no %s app domain in your account", utility.Red(args[1]))
					os.Exit(1)
				}
			}
			appDomainList = append(appDomainList, domain)
		} else {
			for _, v := range args[1:] {
				domain, err := client.FindAppDomain(v, findApp.ID)
				if err == nil {
					appDomainList = append(appDomainList, domain)
				}
			}
		}

		domainNameList := []string{}
		for _, v := range appDomainList {
			domainNameList = append(domainNameList, v)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("app %s", pluralize.Pluralize(len(appDomainList), "domain")), defaultYes, strings.Join(domainNameList, ", ")) {
			for _, v := range appDomainList {
				domain, _ := client.FindAppDomain(v, findApp.ID)
				_, err = client.DeleteAppDomain(appDomainList, findApp.ID, domain)
				if err != nil {
					utility.Error("error deleting the App domain: %s", err)
					os.Exit(1)
				}
			}
			ow := utility.NewOutputWriter()

			for _, v := range appDomainList {
				ow.StartLine()
				ow.AppendDataWithLabel("domain_name", v, "App Domain Name")
			}

			switch outputFormat {
			case "json":
				if len(appDomainList) == 1 {
					ow.WriteSingleObjectJSON(prettySet)
				} else {
					ow.WriteMultipleObjectsJSON(prettySet)
				}
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The domain %s (%s) has been deleted\n", pluralize.Pluralize(len(appDomainList), "domain"), strings.Join(domainNameList, ", "))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
