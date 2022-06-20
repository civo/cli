package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var objectStoreUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit", "modify", "change", "resize"},
	Short:   "Update an Object Store",
	Example: "civo objectstore update OBJECTSTORE_NAME --size SIZE --max-objects MAX_OBJECTS",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()

		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findObjectStore, err := client.FindObjectStore(args[0])
		if err != nil {
			utility.Error("Object Store %s", err)
			os.Exit(1)
		}

		if bucketSize == 0 && maxObjects == 0 {
			utility.Error("You must specify either a size to update your Object Store or max objects to update your Object Store")
			os.Exit(1)
		}

		objectStore, err := client.UpdateObjectStore(findObjectStore.ID, &civogo.UpdateObjectStoreRequest{
			MaxSizeGB:  bucketSize,
			MaxObjects: maxObjects,
			Region:     client.Region,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": objectStore.ID, "name": findObjectStore.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			if bucketSize != 0 && maxObjects == 1000 {
				fmt.Printf("The Object Store with ID %s was updated to size: %s GB\n", utility.Green(objectStore.ID), strconv.Itoa(bucketSize))
				os.Exit(0)
			} else if maxObjects != 0 && bucketSize == 500 {
				fmt.Printf("The Object Store with ID %s was updated to max objects: %s\n", utility.Green(objectStore.ID), strconv.Itoa(maxObjects))
				os.Exit(0)
			}
			fmt.Printf("The Object Store with ID %s was updated to size: %s GB, objects: %s\n", utility.Green(objectStore.ID), strconv.Itoa(bucketSize), strconv.Itoa(maxObjects))
		}
	},
}
