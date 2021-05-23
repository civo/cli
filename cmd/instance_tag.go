package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceTagCmd = &cobra.Command{
	Use:     "tag",
	Example: "civo instance tag ID/HOSTNAME tag1 tag2 tag3",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"tags"},
	Short:   "Change the instance's tags",
	Long: `Change the tags for an instance with partial ID/name provided.
If you wish to use a custom format, the available fields are:

	* id
	* hostname
	* reverse_dns
	* tags`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		tags := strings.Join(args[1:], " ")

		_, err = client.SetInstanceTags(instance, tags)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) has been tagged with '%s'\n", utility.Green(instance.Hostname), instance.ID, utility.Green(tags))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", instance.ID, "ID")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			ow.AppendDataWithLabel("reverse_dns", instance.ReverseDNS, "Reverse DNS")
			ow.AppendDataWithLabel("notes", instance.Notes, "Notes")
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON(prettySet)
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
