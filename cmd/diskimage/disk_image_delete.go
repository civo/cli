package diskimage

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var diskImageDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove"},
	Example: `civo diskimage delete ID/NAME`,
	Short:   "Delete a disk image",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		diskImg, err := client.FindDiskImage(args[0])
		if err != nil {
			utility.Error("Disk Image %s", err)
			os.Exit(1)
		}

		err = client.DeleteDiskImage(diskImg.ID)
		if err != nil {
			utility.Error("Error deleting the disk image: %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.AppendData("Disk image", diskImg.Name)
		ow.AppendData("Result", "deleted")

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
