package kfcluster

import (
	"errors"
	"fmt"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var kfClusterList []utility.ObjecteList
var kfcDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove", "destroy"},
	Short:   "Delete a kubeflow cluster",
	Example: "civo kfc delete <KFCLUSTER-NAME>",
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
			kfc, err := client.FindKfCluster(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s kfcluster in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one kfcluster with that name in your account")
					os.Exit(1)
				}
			}
			kfClusterList = append(kfClusterList, utility.ObjecteList{ID: kfc.ID, Name: kfc.Name})
		} else {
			for _, v := range args {
				kfc, err := client.FindKfCluster(v)
				if err == nil {
					kfClusterList = append(kfClusterList, utility.ObjecteList{ID: kfc.ID, Name: kfc.Name})
				}
			}
		}

		kfcNameList := []string{}
		for _, v := range kfClusterList {
			kfcNameList = append(kfcNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(kfClusterList), "KfCluster"), common.DefaultYes, strings.Join(kfcNameList, ", ")) {

			for _, v := range kfClusterList {
				kfc, err := client.FindKfCluster(v.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				_, err = client.DeleteKfCluster(kfc.ID)
				if err != nil {
					utility.Error("Error deleting the KfCluster: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range kfClusterList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("kfcluster", v.Name, "KfCluster")
			}

			switch common.OutputFormat {
			case "json":
				if len(kfClusterList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(kfClusterList), "kfcluster"), utility.Green(strings.Join(kfcNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
