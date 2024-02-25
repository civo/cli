package diskimage

import (
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var diskImageFindCmd = &cobra.Command{
	Use:     "find",
	Aliases: []string{"get", "search"},
	Example: `civo diskimage find`,
	Short:   "Finds a disk image by either part of the ID or part of the name",
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

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": diskImg.ID, "name": diskImg.Name, "version": diskImg.Version, "state": diskImg.State, "distribution": diskImg.Distribution})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
