package objectstore

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var bucketSize int64
var owner string
var waitOS bool

var objectStoreCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo objectstore create OBJECTSTORE_NAME --size SIZE",
	Short:   "Create a new Object Store",
	Long:    "Bucket size should be in Gigabytes (GB) and must be a multiple of 500, starting from 500.\n",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		var objectStoreName string
		if len(args) > 0 {
			if utility.ValidNameLength(args[0]) {
				utility.Warning("the bucket name cannot be longer than 63 characters")
				os.Exit(1)
			}
			objectStoreName = args[0]

		} else {
			objectStoreName = utility.RandomName()
		}

		check, region, err := utility.CheckAvailability("object_store", common.RegionSet)
		if err != nil {
			utility.Error("Error checking availability %s", err)
			os.Exit(1)
		}

		if !check {
			utility.Error("Sorry you can't create an object store in the %s region", region)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		if bucketSize < 500 {
			utility.Error("The minimum size to create an object store is 500 GB. Please provide a valid size.")
			os.Exit(1)
		} else if bucketSize%500 != 0 {
			utility.YellowConfirm("The size to create an object store must be a multiple of 500. Would you like to create an %s of %d GB instead? (y/n) ? ", utility.Green("object store"), bucketSize+(500-bucketSize%500))
			_, err := utility.UserAccepts(os.Stdin)
			if err != nil {
				utility.Error("Unable to parse the input: %s", err)
				os.Exit(1)
			}
			bucketSize = bucketSize + (500 - bucketSize%500)
		}

		var credential *civogo.ObjectStoreCredential
		var store *civogo.ObjectStore
		if owner != "" {
			credential, err = client.FindObjectStoreCredential(owner)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			store, err = client.NewObjectStore(&civogo.CreateObjectStoreRequest{
				Name:        objectStoreName,
				MaxSizeGB:   bucketSize,
				AccessKeyID: credential.AccessKeyID,
				Region:      client.Region,
			})
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		} else {
			store, err = client.NewObjectStore(&civogo.CreateObjectStoreRequest{
				Name:      objectStoreName,
				MaxSizeGB: bucketSize,
				Region:    client.Region,
			})
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
		}

		var executionTime string
		if waitOS {
			startTime := utility.StartTime()
			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Writer = os.Stderr
			s.Prefix = fmt.Sprintf("Creating an Object Store with maxSize %d, called %s... ", store.MaxSize, store.Name)
			s.Start()

			for stillCreating {
				storeCheck, err := client.FindObjectStore(store.ID)
				if err != nil {
					utility.Error("Object Store %s", err)
					os.Exit(1)
				}
				if storeCheck.Status == "ready" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}

			executionTime = utility.TrackTime(startTime)
		}

		objectStore, err := client.FindObjectStore(objectStoreName)
		if err != nil {
			utility.Error("ObjectStore %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"name": objectStore.Name, "id": objectStore.ID, "access_key": objectStore.OwnerInfo.AccessKeyID})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if waitOS {
				fmt.Printf("Created Object Store %s in %s in %s\n", utility.Green(objectStore.Name), utility.Green(client.Region), executionTime)
				fmt.Printf("Created default admin credentials, access key is %s, this will be deleted if the Object Store is deleted. ", utility.Green(objectStore.OwnerInfo.AccessKeyID))
				fmt.Printf("To access the secret key run: civo objectstore credential secret --access-key=%s\n", utility.Green(objectStore.OwnerInfo.AccessKeyID))
			} else {
				fmt.Printf("Creating Object Store %s in %s\n", utility.Green(objectStore.Name), utility.Green(client.Region))
				fmt.Printf("To check the status of the Object Store run: civo objectstore show %s\n", utility.Green(objectStore.Name))
			}
		}
	},
}
