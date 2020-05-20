package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var templateRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo template remove CODE",
	Short:   "Remove a template",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("template") == nil {
			template, err := client.GetTemplateByCode(args[0])
			if err != nil {
				utility.Error("Finding the template for your search failed with %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteTemplate(template.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": template.ID, "Name": template.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The template called %s with ID %s was deleted\n", utility.Green(template.Name), utility.Green(template.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
