package cmd

import (
	"fmt"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var configName, configValue string

var appConfigSetCmd = &cobra.Command{
	Use:     "set",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Set application config",
	Example: "civo app config set APP_NAME foo=bar foo2=bar2",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findApp, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if configName == "PORT" {
			utility.Error("PORT is an immutable field picked by Civo.")
			os.Exit(1)
		}

		updatedConfig := make([]civogo.EnvVar, len(findApp.Config))
		for _, arg := range args[1:] {
			if strings.Contains(arg, "=") {
				parts := strings.Split(arg, "=")
				if len(parts) != 2 {
					utility.Error("Invalid argument %s", arg)
					os.Exit(1)
				}

				updatedConfig = append(updatedConfig, civogo.EnvVar{
					Name:  parts[0],
					Value: parts[1],
				})
			}
			for _, envVar := range findApp.Config {
				for _, newEnvVar := range updatedConfig {
					if envVar.Name == newEnvVar.Name {
						envVar.Value = newEnvVar.Value
					}
				}
			}
		}

		config := &civogo.UpdateApplicationRequest{
			Config: updatedConfig,
		}

		app, err := client.UpdateApplication(findApp.ID, config)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": app.ID, "name": app.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Application %s's config has been updated.\n", utility.Green(app.Name))
		}
	},
}
