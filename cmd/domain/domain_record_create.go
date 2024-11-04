package domain

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var recordName, recordType, recordValue string
var recordTTL, recordPriority int

var domainRecordCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new domain record",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo domain record create DOMAIN/DOMAIN_ID [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		domain, err := client.FindDNSDomain(args[0])
		if err != nil {
			utility.Error("Unable to find the domain for your search %s", err)
			os.Exit(1)
		}

		newRecordConfig := &civogo.DNSRecordConfig{
			Name:     recordName,
			Value:    recordValue,
			TTL:      recordTTL,
			Priority: recordPriority,
		}

		// Sanitace the record type
		recordType = strings.ReplaceAll(recordType, " ", "")

		if recordType == "A" || recordType == "a" || recordType == "alias" {
			newRecordConfig.Type = civogo.DNSRecordTypeA
		}

		if recordType == "CNAME" || recordType == "cname" || recordType == "canonical" {
			newRecordConfig.Type = civogo.DNSRecordTypeCName
		}

		if recordType == "MX" || recordType == "mx" || recordType == "mail" {
			newRecordConfig.Type = civogo.DNSRecordTypeMX
		}

		if recordType == "TXT" || recordType == "txt" || recordType == "text" {
			newRecordConfig.Type = civogo.DNSRecordTypeTXT
		}

		if recordType == "SRV" || recordType == "srv" || recordType == "service" {
			newRecordConfig.Type = civogo.DNSRecordTypeSRV
		}

		if recordType == "NS" || recordType == "ns" || recordType == "nameserver" {
			newRecordConfig.Type = civogo.DNSRecordTypeNS
		}

		record, err := client.CreateDNSRecord(domain.ID, newRecordConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": domain.ID, "name": domain.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Created %s record %s for %s with a TTL of %s seconds and with a priority of %s with ID %s", utility.Green(string(record.Type)), utility.Green(record.Name), utility.Green(domain.Name), utility.Green(strconv.Itoa(record.TTL)), utility.Green(strconv.Itoa(record.Priority)), utility.Green(record.ID))
		}
	},
}
