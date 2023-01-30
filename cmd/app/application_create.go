package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appSize, gitURL, image, tagName, branchName string
var wait bool

var appCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo app create APP_NAME [flags]",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Create a new application",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		config, err := client.NewApplicationConfig()
		if err != nil {
			utility.Error("Unable to create a new config for the app %s", err)
			os.Exit(1)
		}

		config.PublicIPv4Required = true

		if gitURL != "" {
			config.GitInfo.GitURL = gitURL
			config.GitInfo.PullPreference.Branch = &branchName
			config.GitInfo.PullPreference.Tag = &tagName
			if os.Getenv("GIT_TOKEN") == "" {
				utility.Error("GIT_TOKEN env var not found %s", err)
				os.Exit(1)
			}
			config.GitInfo.GitToken = os.Getenv("GIT_TOKEN")
		}

		if image != "" {
			config.Image = &image
		}

		if config.Image == nil && config.GitInfo == nil {
			utility.Error("No image or git info specified for the app %s", err)
			os.Exit(1)
		}

		if utility.ValidNameLength(args[0]) {
			utility.Warning("the name cannot be longer than 40 characters")
			os.Exit(1)
		}
		config.Name = args[0]

		// if len(args) > 0 {
		// 	if utility.ValidNameLength(args[0]) {
		// 		utility.Warning("the name cannot be longer than 40 characters")
		// 		os.Exit(1)
		// 	}
		// 	config.Name = args[0]
		// }

		if appSize != "" {
			config.Size = appSize
		} else {
			config.Size = "small"
		}

		var executionTime string
		startTime := utility.StartTime()

		var application *civogo.Application
		resp, err := client.CreateApplication(config)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if wait {
			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating application (%s)... ", resp.Name)
			s.Start()

			for stillCreating {
				application, err = client.FindApplication(resp.Name)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				if strings.ToLower(application.Status) == "ready" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
			executionTime = utility.TrackTime(startTime)
		} else {
			application, err = client.FindApplication(resp.Name)
			if err != nil {
				utility.Error("App %s", err)
				os.Exit(1)
			}
		}

		if common.OutputFormat == "human" {
			if executionTime != "" {
				fmt.Printf("The app %s has been created in %s\n", utility.Green(application.Name), executionTime)
			} else {
				fmt.Printf("The app %s has been created\n", utility.Green(application.Name))
			}
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", resp.ID, "ID")
			ow.AppendDataWithLabel("name", resp.Name, "Name")
			ow.AppendDataWithLabel("network_id", resp.NetworkID, "Network ID")
			ow.AppendDataWithLabel("size", resp.Size, "Size")
			ow.AppendDataWithLabel("status", resp.Status, "Status")

			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		}
	},
}
