package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

func scaleCmd(cmd *cobra.Command, args []string) {

	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	findApp, err := client.FindApplication(args[0])
	if err != nil {
		utility.Error("App %s", err)
		os.Exit(1)
	}

	processInfo := make([]civogo.ProcessInfo, 0)
	for _, arg := range args[1:] {
		if strings.Contains(arg, "=") {
			parts := strings.Split(arg, "=")
			if len(parts) != 2 {
				utility.Error("Invalid argument %s", arg)
				os.Exit(1)
			}

			procCount, err := strconv.Atoi(parts[1])
			if err != nil {
				utility.Error("Invalid count %s", arg)
				os.Exit(1)
			}

			processInfo = append(processInfo, civogo.ProcessInfo{
				ProcessType:  parts[0],
				ProcessCount: procCount,
			})
		}
	}
	for _, process := range findApp.ProcessInfo {
		for _, newProcess := range processInfo {
			if process.ProcessType == newProcess.ProcessType {
				process.ProcessCount = newProcess.ProcessCount
			}
		}
	}

	application := &civogo.UpdateApplicationRequest{
		ProcessInfo: processInfo,
	}

	app, err := client.UpdateApplication(findApp.ID, application)
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
		fmt.Printf("The application %s has been updated.\n", utility.Green(app.Name))
	}
}
