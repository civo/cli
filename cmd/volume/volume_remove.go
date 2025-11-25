package volume

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/pkg/pluralize"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var volumeRemoveCmdExamples = []string{
	"civo volume rm VOLUME_NAME",
	"civo volume rm VOLUME_ID",
}

var volumesList []utility.ObjecteList
var volumeRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: strings.Join(volumeRemoveCmdExamples, "\n"),
	Short:   "Remove a volume",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		if len(args) == 1 {
			volume, err := client.FindVolume(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s volume in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one volume with that value in your account")
					os.Exit(1)
				}
			}
			volumesList = append(volumesList, utility.ObjecteList{ID: volume.ID, Name: volume.Name})
		} else {
			for _, v := range args {
				volume, err := client.FindVolume(v)
				if err == nil {
					volumesList = append(volumesList, utility.ObjecteList{ID: volume.ID, Name: volume.Name})
				}
			}
		}

		volumeNameList := []string{}
		for _, v := range volumesList {
			volumeNameList = append(volumeNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(volumesList), "Volume"), common.DefaultYes, strings.Join(volumeNameList, ", ")) {

			for _, v := range volumesList {
				vol, err := client.FindVolume(v.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				_, err = client.DeleteVolume(vol.ID)
				if err != nil {
					utility.Error("Error deleting the Volume: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, volume := range volumesList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", volume.ID, "ID")
				ow.AppendDataWithLabel("volume", volume.Name, "Volume")
			}

			switch common.OutputFormat {
			case "json":
				if len(volumesList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) %s been deleted\n",
					pluralize.Pluralize(len(volumesList), "volume"),
					utility.Green(strings.Join(volumeNameList, ", ")),
					pluralize.Has(len(volumeNameList)))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
