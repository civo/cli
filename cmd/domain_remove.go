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

var domainList []utility.ObjecteList
var domainRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Short:   "Remove a domain",
	Example: "civo domain remove DOMAIN/DOMAIN_ID",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			domain, err := client.FindDNSDomain(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s domain in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one domain with that name in your account")
					os.Exit(1)
				}
			}
			domainList = append(domainList, utility.ObjecteList{ID: domain.ID, Name: domain.Name})
		} else {
			for _, v := range args {
				domain, err := client.FindDNSDomain(v)
				if err == nil {
					domainList = append(domainList, utility.ObjecteList{ID: domain.ID, Name: domain.Name})
				}
			}
		}

		domainNameList := []string{}
		for _, v := range domainList {
			domainNameList = append(domainNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(domainList), "domain"), defaultYes, strings.Join(domainNameList, ", ")) {

			for _, v := range domainList {
				domain, _ := client.FindDNSDomain(v.ID)
				_, err = client.DeleteDNSDomain(domain)
				if err != nil {
					utility.Error("error deleting the DNS domain: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range domainList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("domain", v.Name, "Domain")
			}

			switch outputFormat {
			case "json":
				if len(domainList) == 1 {
					ow.WriteSingleObjectJSON(prettySet)
				} else {
					ow.WriteMultipleObjectsJSON(prettySet)
				}
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(domainList), "domain"), utility.Green(strings.Join(domainNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
