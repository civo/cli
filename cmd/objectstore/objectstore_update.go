package objectstore

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var objectStoreUpdateCmd = &cobra.Command{
	Use:     "resize",
	Aliases: []string{"edit", "modify", "change", "update"},
	Short:   "Update an Object Store",
	Example: "civo objectstore resize OBJECTSTORE_NAME --size SIZE",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		findObjectStore, err := client.FindObjectStore(args[0])
		if err != nil {
			utility.Error("Object Store %s", err)
			os.Exit(1)
		}

		if bucketSize == 0 {
			utility.Error("You must specify size to update your Object Store")
			os.Exit(1)
		}

		var credential *civogo.ObjectStoreCredential
		if owner != "" {
			credential, err = client.FindObjectStoreCredential(owner)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		objectStore, err := client.UpdateObjectStore(findObjectStore.ID, &civogo.UpdateObjectStoreRequest{
			MaxSizeGB:   bucketSize,
			AccessKeyID: credential.AccessKeyID,
			Region:      client.Region,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": objectStore.ID, "name": findObjectStore.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if bucketSize != 0 && owner != "" {
				fmt.Printf("The Object Store with ID %s was updated to size: %d GB and owner %s \n", utility.Green(objectStore.ID), bucketSize, credential.Name)
				os.Exit(0)
			} else if bucketSize != 0 && owner == "" {
				fmt.Printf("The Object Store with ID %s was updated to size: %d GB \n", utility.Green(objectStore.ID), bucketSize)
				os.Exit(0)
			} else {
				fmt.Printf("The owner of Object Store with ID %s was updated to:%s \n", utility.Green(objectStore.ID), credential.Name)
			}
		}
	},
}
