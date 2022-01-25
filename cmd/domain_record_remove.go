package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var domainRecordList []utility.ObjectList
var domainRecordRemoveCmd = &cobra.Command{
	Use:     "remove [DOMAIN|DOMAIN_ID] [RECORD_ID]",
	Aliases: []string{"delete", "destroy", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove domain record",
	Example: "civo domain record remove DOMAIN/DOMAIN_ID RECORD_ID",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllDomainList(), cobra.ShellCompDirectiveNoFileComp
		}
		return getDomainList(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

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

		if len(args) == 2 {
			record, err := client.GetDNSRecord(domain.ID, args[1])
			if err != nil {
				if errors.Is(err, civogo.ErrDNSRecordNotFound) {
					utility.Error("sorry there is no %s domain record in your account", utility.Red(args[1]))
					os.Exit(1)
				}
			}
			domainRecordList = append(domainRecordList, utility.ObjectList{ID: record.ID, Name: record.Name})
		} else {
			for _, v := range args[1:] {
				record, err := client.GetDNSRecord(domain.ID, v)
				if err == nil {
					domainRecordList = append(domainRecordList, utility.ObjectList{ID: record.ID, Name: record.Name})
				}
			}
		}

		domainRecordNameList := []string{}
		for _, v := range domainRecordList {
			domainRecordNameList = append(domainRecordNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("domain %s", pluralize.Pluralize(len(domainRecordList), "record")), defaultYes, strings.Join(domainRecordNameList, ", ")) {

			for _, v := range domainRecordList {
				record, _ := client.GetDNSRecord(domain.ID, v.ID)
				_, err = client.DeleteDNSRecord(record)
				if err != nil {
					utility.Error("error deleting the DNS record: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range domainRecordList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch outputFormat {
			case "json":
				if len(domainRecordList) == 1 {
					ow.WriteSingleObjectJSON(prettySet)
				} else {
					ow.WriteMultipleObjectsJSON(prettySet)
				}
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The domain %s (%s) has been deleted\n", pluralize.Pluralize(len(domainRecordList), "record"), strings.Join(domainRecordNameList, ", "))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
