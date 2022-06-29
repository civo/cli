package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceUpgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Example: "civo instance upgrade ID/HOSTNAME g3.xlarge",
	Args:    cobra.MinimumNArgs(2),
	Aliases: []string{"resize"},
	Short:   "Upgrade an instance",
	Long: `Upgrade instance with ID to size provided. Downgrades to smaller sizes are not possible.
Run civo sizes for all the size names.
If you wish to use a custom format, the available fields are:

	* id
	* hostname
	* old_size
	* new_size`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		if len(args) != 2 {
			fmt.Printf("You must specify %s parameters (you gave %s), the ID/name and the new size\n", utility.Red("2"), utility.Red(strconv.Itoa(len(args))))
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
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

		sizes, err := client.ListInstanceSizes()
		if err != nil {
			utility.Error("Checking size is valid failed with %s", err)
			os.Exit(1)
		}

		var resizing bool
		for _, size := range sizes {
			if size.Name == args[1] {
				resizing = true
				_, err = client.UpgradeInstance(instance.ID, size.Name)
				if err != nil {
					utility.Error("Upgrading instance failed with %s", err)
					os.Exit(1)
				}
			}
		}

		if !resizing {
			utility.Error("Unable to find size %s", args[1])
			os.Exit(1)
		}

		if common.OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being upgraded to %s\n", utility.Green(instance.Hostname), instance.ID, utility.Green(args[1]))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", instance.ID, "ID")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			ow.AppendDataWithLabel("old_size", instance.Size, "Old Size")
			ow.AppendDataWithLabel("new_size", args[1], "New Size")
			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		}
	},
}
