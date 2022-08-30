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

var credentialSize int

var objectStoreCredentialCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new Object Store Credential",
	Example: "civo objectstore credential create CREDENTIAL_NAME --size SIZE",
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

		credential, err := client.NewObjectStoreCredential(&civogo.CreateObjectStoreCredentialRequest{
			Name:      args[0],
			MaxSizeGB: &credentialSize,
			//Region:    client.Region,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		var executionTime string
		if waitOS {
			startTime := utility.StartTime()
			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating an Object Store Credential with maxSize %d, called %s... ", credential.MaxSizeGB, credential.Name)
			s.Start()

			for stillCreating {
				credCheck, err := client.FindObjectStoreCredential(credential.ID)
				if err != nil {
					utility.Error("Object Store Credential %s", err)
					os.Exit(1)
				}
				if credCheck.Status == "ready" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
			executionTime = utility.TrackTime(startTime)
		}

		objectStoreCred, err := client.FindObjectStoreCredential(args[0])
		if err != nil {
			utility.Error("ObjectStore Credential %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"name": objectStoreCred.Name, "id": objectStoreCred.ID})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if waitOS {
				fmt.Printf("Created Object Store Credential %s in %s in %s\n", utility.Green(objectStoreCred.Name), utility.Green(client.Region), executionTime)
			} else {
				fmt.Printf("Creating Object Store Credential %s in %s\n", utility.Green(objectStoreCred.Name), utility.Green(client.Region))
			}
		}
	},
}
