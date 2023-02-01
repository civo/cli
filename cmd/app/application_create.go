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
	Short:   "Create a new application",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		// TODO: Add network and firewall flags

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

		if image == "" && gitURL == "" {
			utility.Error("No image or git info specified for the app")
			os.Exit(1)
		}

		if gitURL != "" {
			// TODO : Add in help menu
			if os.Getenv("GIT_TOKEN") == "" {
				utility.Error("GIT_TOKEN env var not found")
				os.Exit(1)
			}
			config.GitInfo = &civogo.GitInfo{}
			config.GitInfo.GitToken = os.Getenv("GIT_TOKEN")
			config.GitInfo.GitURL = gitURL
			pullPref := &civogo.PullPreference{}
			if branchName != "" {
				pullPref.Branch = &branchName
			}
			if tagName != "" {
				pullPref.Tag = &tagName
			}
			config.GitInfo.PullPreference = pullPref
		}

		if image != "" {
			config.Image = &image
		}

		if len(args) > 0 {
			if args[0] != "" {
				if utility.ValidNameLength(args[0]) {
					utility.Warning("the name cannot be longer than 40 characters")
					os.Exit(1)
				}
				config.Name = args[0]
			}
		}

		if appSize != "" {
			config.Size = appSize
		}

		var application *civogo.Application
		resp, err := client.CreateApplication(config)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		var executionTime string
		if wait {
			startTime := utility.StartTime()
			stillCreating := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = fmt.Sprintf("Creating application (%s)... ", resp.Name)
			s.Start()

			var statusAvailable bool
			for stillCreating {
				application, err = client.FindApplication(resp.Name)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				if strings.ToLower(application.Status) == "available" && !statusAvailable {
					statusAvailable = true
					fmt.Println("Environment has been set up, deploying workloads now.")
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
				fmt.Printf("The app %s is deploying\n", utility.Green(application.Name))
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
