package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var domainRecordRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "destroy", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove record",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("domain record") == nil {
			domain, err := client.FindDNSDomain(args[0])
			if err != nil {
				fmt.Printf("Unable to find domain for your search: %s\n", aurora.Red(err))
				os.Exit(1)
			}

			record, err := client.GetDNSRecord(domain.ID, args[1])
			if err != nil {
				fmt.Printf("Unable to get domains record: %s\n", aurora.Red(err))
				os.Exit(1)
			}

			_, err = client.DeleteDNSRecord(record)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": record.ID, "Name": record.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The domain record called %s with ID %s was delete\n", aurora.Green(record.Name), aurora.Green(record.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
