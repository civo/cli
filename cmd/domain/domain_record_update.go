package domain

import (
	"errors"
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

var updateRecordName, updateRecordType, updateRecordValue string
var updateRecordTTL, updateRecordPriority int

var domainRecordUpdateCmd = &cobra.Command{
	Use:     "update [DOMAIN|DOMAIN_ID] [RECORD_ID]",
	Aliases: []string{"change", "modify"},
	Short:   "Update a domain record",
	Args:    cobra.MinimumNArgs(2),
	Example: "civo domain record update DOMAIN/DOMAIN_ID RECORD_ID [flags]",
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
			utility.Error("Unable to find the domain for your search %s", err)
			os.Exit(1)
		}

		record, err := client.GetDNSRecord(domain.ID, args[1])
		if err != nil {
			if errors.Is(err, civogo.ErrDNSRecordNotFound) {
				utility.Error("sorry there is no %s domain record in your account", utility.Red(args[1]))
				os.Exit(1)
			}
			utility.Error("%s", err)
			os.Exit(1)
		}

		// Seed the config from the existing record so unspecified fields are preserved
		// (the API call is a full-replace PUT).
		recordConfig := &civogo.DNSRecordConfig{
			Type:     record.Type,
			Name:     record.Name,
			Value:    record.Value,
			TTL:      record.TTL,
			Priority: record.Priority,
		}

		if cmd.Flags().Changed("name") {
			recordConfig.Name = updateRecordName
		}

		if cmd.Flags().Changed("value") {
			recordConfig.Value = updateRecordValue
		}

		if cmd.Flags().Changed("ttl") {
			recordConfig.TTL = updateRecordTTL
		}

		if cmd.Flags().Changed("priority") {
			recordConfig.Priority = updateRecordPriority
		}

		if cmd.Flags().Changed("type") {
			// Sanitise the record type
			updateRecordType = strings.ReplaceAll(updateRecordType, " ", "")

			if updateRecordType == "A" || updateRecordType == "a" || updateRecordType == "alias" {
				recordConfig.Type = civogo.DNSRecordTypeA
			}

			if updateRecordType == "CNAME" || updateRecordType == "cname" || updateRecordType == "canonical" {
				recordConfig.Type = civogo.DNSRecordTypeCName
			}

			if updateRecordType == "MX" || updateRecordType == "mx" || updateRecordType == "mail" {
				recordConfig.Type = civogo.DNSRecordTypeMX
			}

			if updateRecordType == "TXT" || updateRecordType == "txt" || updateRecordType == "text" {
				recordConfig.Type = civogo.DNSRecordTypeTXT
			}

			if updateRecordType == "SRV" || updateRecordType == "srv" || updateRecordType == "service" {
				recordConfig.Type = civogo.DNSRecordTypeSRV
			}

			if updateRecordType == "NS" || updateRecordType == "ns" || updateRecordType == "nameserver" {
				recordConfig.Type = civogo.DNSRecordTypeNS
			}
		}

		updatedRecord, err := client.UpdateDNSRecord(record, recordConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := dnsRecordOutputWriter(updatedRecord)

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Updated %s record %s for %s with a TTL of %s seconds and with a priority of %s with ID %s\n", utility.Green(strings.ToUpper(string(updatedRecord.Type))), utility.Green(updatedRecord.Name), utility.Green(domain.Name), utility.Green(strconv.Itoa(updatedRecord.TTL)), utility.Green(strconv.Itoa(updatedRecord.Priority)), utility.Green(updatedRecord.ID))
		}
	},
}
