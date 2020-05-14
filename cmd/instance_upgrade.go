package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceUpgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Example: "civo instance upgrade ID/HOSTNAME g2.xlarge",
	Aliases: []string{"resize"},
	Short:   "Upgrade an instance",
	Long: `Upgrade instance with ID to size provided. Downgrades to smaller sizes are not possible.
Run civo sizes for all the size names.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* OldSize
	* NewSize`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("You must specify %s parameters (you gave %s), the ID/name and the new size\n", utility.Red("2"), utility.Red(strconv.Itoa(len(args))))
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Finding instance %s", err)
			os.Exit(1)
		}

		sizes, err := client.ListInstanceSizes()
		if err != nil {
			utility.Error("Checking size %s", err)
			os.Exit(1)
		}

		var resizing bool
		for _, size := range sizes {
			if size.Name == args[1] {
				resizing = true
				_, err = client.UpgradeInstance(instance.ID, size.Name)
				if err != nil {
					utility.Error("Upgrading instance %s", err)
					os.Exit(1)
				}
			}
		}

		if !resizing {
			utility.Error("Unable to find size", args[1])
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being upgraded to %s\n", utility.Green(instance.Hostname), instance.ID, utility.Green(args[1]))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendDataWithLabel("OldSize", instance.Size, "Old Size")
			ow.AppendDataWithLabel("NewSize", args[1], "New Size")
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
