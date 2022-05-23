package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appName, appSize, appSSHKeyIDs string

var appCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo app create APP_NAME [flags]",
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

		if appName != "" {
			if utility.ValidNameLength(appName) {
				utility.Warning("the name cannot be longer than 63 characters")
				os.Exit(1)
			}
			config.Name = appName
		}

		if appSSHKeyIDs != "" {
			config.SSHKeyIDs = strings.Split(appSSHKeyIDs, ",")
		}

		if len(args) > 0 {
			if utility.ValidNameLength(args[0]) {
				utility.Warning("the name cannot be longer than 63 characters")
				os.Exit(1)
			}
			config.Name = args[0]
		}

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
				if strings.ToLower(application.Status) == "available" {
					stillCreating = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
			executionTime = utility.TrackTime(startTime)
		} else {
			// we look for the created app to obtain the data that we need
			application, err = client.FindApplication(resp.Name)
			if err != nil {
				utility.Error("App %s", err)
				os.Exit(1)
			}
		}

		if outputFormat == "human" {
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
			ow.AppendDataWithLabel("description", resp.Description, "Description")
			//ow.AppendDataWithLabel("image", resp.Image, "Image")
			ow.AppendDataWithLabel("size", resp.Size, "Size")
			ow.AppendDataWithLabel("status", resp.Status, "Status")
			// ow.AppendDataWithLabel("process_info", resp.ProcessInfo, "Process Info")
			ow.AppendDataWithLabel("domains", strings.Join(resp.Domains, ", "), "Domains")
			ow.AppendDataWithLabel("ssh_key_ids", strings.Join(resp.SSHKeyIDs, ", "), "SSH Key IDs")
			//ow.AppendDataWithLabel("config", resp.Config, "Config")

			if outputFormat == "json" {
				ow.WriteSingleObjectJSON(prettySet)
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
