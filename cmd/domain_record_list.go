package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var domainRecordListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "List all domains records",
	Long: `List all current domain records.
If you wish to use a custom format, the available fields ar	:

	* ID
	* Name
	* Value
	* Type
	* TTL
	* Priority	

Example: civo domain record ls DOMAIN or DOMAIN_ID -o custom -f "ID: Name"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		domain, err := client.FindDNSDomain(args[0])
		if err != nil {
			utility.Error("Unable to find domain for your search %s", err)
			os.Exit(1)
		}

		records, err := client.ListDNSRecords(domain.ID)
		if err != nil {
			utility.Error("Unable to list domains %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, record := range records {
			ow.StartLine()

			ow.AppendData("ID", record.ID)
			ow.AppendData("Name", record.Name)
			ow.AppendData("Value", record.Value)
			ow.AppendData("Type", string(record.Type))
			ow.AppendData("TTL", strconv.Itoa(record.TTL))
			ow.AppendData("Priority", strconv.Itoa(record.Priority))

		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
