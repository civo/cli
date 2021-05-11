package cmd

import (
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var domainRecordShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Show domain record",
	Example: "civo domain record show DOMAIN/DOMAIN_ID RECORD_ID",
	Long: `Show the specified record.
If you wish to use a custom format, the available fields are:

	* ID
	* DomainID
	* Name
	* Value
	* Type
	* TTL
	* Priority

Example: civo domain record show RECORD_ID -o custom -f "ID: Name"`,
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

		record, err := client.GetDNSRecord(domain.ID, args[1])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()

		ow.AppendData("ID", record.ID)
		ow.AppendDataWithLabel("DomainID", record.DNSDomainID, "Domain ID")
		ow.AppendData("Name", record.Name)
		ow.AppendData("Value", record.Value)

		if record.Type == "a" {
			ow.AppendData("Type", string(civogo.DNSRecordTypeA))
		}

		if record.Type == "cname" {
			ow.AppendData("Type", string(civogo.DNSRecordTypeCName))
		}

		if record.Type == "mx" {
			ow.AppendData("Type", string(civogo.DNSRecordTypeMX))
		}

		if record.Type == "txt" {
			ow.AppendData("Type", string(civogo.DNSRecordTypeTXT))
		}

		if record.Type == "srv" {
			ow.AppendData("Type", string(civogo.DNSRecordTypeSRV))
		}

		ow.AppendData("TTL", strconv.Itoa(record.TTL))
		ow.AppendData("Priority", strconv.Itoa(record.Priority))
		ow.AppendDataWithLabel("CreatedAt", record.CreatedAt.Format(time.RFC1123), "Created At")
		ow.AppendDataWithLabel("UpdatedAt", record.UpdatedAt.Format(time.RFC1123), "Updated At")

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
