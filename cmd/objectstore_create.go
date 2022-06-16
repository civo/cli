package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var maxObjects int
var bucketSize int

var objectStoreCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo objectstore create OBJECTSTORE_NAME_PREFIX --size SIZE",
	Short:   "Create a new Object Store",
	Long:    "Bucket Size should be in Gigabytes (GB) and must be a multiple of 500, starting from 500.\n An objectstore name will be generated from the prefix provided.\n",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if bucketSize == 0 {
			bucketSize = 500
		}
		_, err = client.NewObjectStore(&civogo.CreateObjectStorageRequest{
			Name:       args[0],
			MaxSizeGB:  bucketSize,
			MaxObjects: maxObjects,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		objectStore, err := client.FindObjectStore(args[0])
		if err != nil {
			utility.Error("ObjectStore %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"name": objectStore.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created an Objectstore with Name %s\n", utility.Green(objectStore.Name))
		}
	},
}
