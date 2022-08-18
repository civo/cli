package instance

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

// InstanceList is a tmp list to hold all instance to delete

var instanceList []utility.ObjecteList
var instanceRemoveCmd = &cobra.Command{
	Use:     "remove",
	Example: "civo instance remove ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"delete", "destroy", "rm"},
	Short:   "Remove/delete instance",
	Long: `Remove the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* id
	* hostname`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			instance, err := client.FindInstance(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s instance in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one instance with that name in your account")
					os.Exit(1)
				}
			}

			instanceList = append(instanceList, utility.ObjecteList{ID: instance.ID, Name: instance.Hostname})

		} else {
			for _, v := range args {
				instance, err := client.FindInstance(v)
				if err == nil {
					instanceList = append(instanceList, utility.ObjecteList{ID: instance.ID, Name: instance.Hostname})
				}
			}
		}

		instanceNameList := []string{}
		for _, v := range instanceList {
			instanceNameList = append(instanceNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(instanceList), "instance"), common.DefaultYes, strings.Join(instanceNameList, ", ")) {

			for _, v := range instanceList {
				_, err = client.DeleteInstance(v.ID)
				if err != nil {
					utility.Error("Error deleting the instance: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range instanceList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("hostname", v.Name, "Hostname")
			}

			switch common.OutputFormat {
			case "json":
				if len(instanceList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(instanceList), "instance"), utility.Green(strings.Join(instanceNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
