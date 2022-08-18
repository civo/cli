package domain

import (
	"os"
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var domainRecordListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: "civo domain record ls DOMAIN/DOMAIN_ID",
	Args:    cobra.MinimumNArgs(1),
	Short:   "List all domains records",
	Long: `List all current domain records.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* value
	* type
	* ttl
	* priority`,
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
			utility.Error("%s", err)
			os.Exit(1)
		}

		records, err := client.ListDNSRecords(domain.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, record := range records {
			ow.StartLine()

			ow.AppendDataWithLabel("id", record.ID, "ID")
			ow.AppendDataWithLabel("name", record.Name, "Name")
			ow.AppendDataWithLabel("value", record.Value, "Value")
			ow.AppendDataWithLabel("type", string(record.Type), "Type")
			ow.AppendDataWithLabel("ttl", strconv.Itoa(record.TTL), "TTL")
			ow.AppendDataWithLabel("priority", strconv.Itoa(record.Priority), "Priority")

		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}

func getDomainList(value string) []string {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	domain, err := client.FindDNSDomain(value)
	if err != nil {
		utility.Error("Unable to list domains %s", err)
		os.Exit(1)
	}

	var domainList []string
	domainList = append(domainList, domain.Name)

	return domainList

}

func getAllDomainList() []string {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	domain, err := client.ListDNSDomains()
	if err != nil {
		utility.Error("Unable to list domains %s", err)
		os.Exit(1)
	}

	var domainList []string
	for _, v := range domain {
		domainList = append(domainList, v.Name)
	}

	return domainList

}
