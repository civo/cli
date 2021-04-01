package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var domainCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new domain",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo domain create NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		domain, err := client.CreateDNSDomain(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": domain.ID, "Name": domain.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a domain called %s with ID %s\n", utility.Green(domain.Name), utility.Green(domain.ID))
			fmt.Println("Please point your domain registrar to Civo nameservers:")
			fmt.Printf("%s\n", utility.Green("ns0.civo.com"))
			fmt.Printf("%s\n", utility.Green("ns1.civo.com"))
		}
	},
}
