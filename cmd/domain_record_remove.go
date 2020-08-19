package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var domainRecordRemoveCmd = &cobra.Command{
	Use:     "remove [DOMAIN|DOMAIN_ID] [RECORD_ID]",
	Aliases: []string{"delete", "destroy", "rm"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Remove domain record",
	Example: "civo domain record remove DOMAIN/DOMAIN_ID RECORD_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		domain, err := client.FindDNSDomain(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry this domain (%s) does not exist in your account", args[0])
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one domain with that name in your account", args[0])
				os.Exit(1)
			}
		}

		record, err := client.GetDNSRecord(domain.ID, args[1])
		if err != nil {
			if errors.Is(err, civogo.ErrDNSRecordNotFound) {
				utility.Error("sorry this domain record (%s) does not exist in your account", args[1])
				os.Exit(1)
			}
		}

		if utility.UserConfirmedDeletion("domain record", defaultYes) == true {

			_, err = client.DeleteDNSRecord(record)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": record.ID, "Name": record.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The domain record called %s with ID %s was delete\n", utility.Green(record.Name), utility.Green(record.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
