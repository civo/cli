package domain

import (
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

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

	* id
	* domain_id
	* name
	* value
	* type
	* ttl
	* priority
	* created_at
	* updated_at

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

		ow := dnsRecordOutputWriter(record)

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
