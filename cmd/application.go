package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:     "app",
	Aliases: []string{"apps, application, applications"},
	Short:   "Manage Applications inside your Civo account",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var appDomainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"domains"},
	Short:   "Details of your application domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var appConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Configure your application",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var appScaleCmd = &cobra.Command{
	Use:     "scale",
	Aliases: []string{"change", "modify", "upgrade"},
	Example: "civo app scale APP-NAME PROCESS-NAME=PROCESS-COUNT",
	Short:   "Scale processes of your application",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

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
		parts := make([]string, 0)
		for _, arg := range args[1:] {
			if strings.Contains(arg, "=") {
				parts := strings.Split(arg, "=")
				if len(parts) != 2 {
					utility.Error("Invalid argument %s", arg)
					os.Exit(1)
				}
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
		for _, process := range findApp.ProcessInfo {
			for _, newProcess := range processInfo {
				if process.ProcessType == newProcess.ProcessType {
					process.ProcessCount = newProcess.ProcessCount
				}
			}
		}

		application := &civogo.UpdateApplicationRequest{
			ProcessInfo: findApp.ProcessInfo,
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
	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appCreateCmd)
	appCreateCmd.Flags().StringVarP(&appName, "name", "n", "", "Name of the application")
	appCreateCmd.Flags().StringVarP(&appSize, "size", "s", "", "Size of the application")
	appCmd.AddCommand(appRemoveCmd)
	appCmd.AddCommand(appScaleCmd)
	appCmd.AddCommand(appRemoteCmd)
	appRemoteCmd.Flags().StringVarP(&remoteName, "remote-name", "r", "", "The name of remote you want to add. E.g. civo")

	//App domain commands
	appCmd.AddCommand(appDomainCmd)
	appDomainCmd.AddCommand(appDomainListCmd)
	appDomainCmd.AddCommand(appDomainAddCmd)
	appDomainCmd.AddCommand(appDomainRemoveCmd)

	//App config commands
	appCmd.AddCommand(appConfigCmd)
	appConfigCmd.AddCommand(appConfigShowCmd)
	appConfigCmd.AddCommand(appConfigSetCmd)
	appConfigSetCmd.Flags().StringVarP(&configName, "name", "n", "", "The name of the environment variable you want to set.")
	appConfigSetCmd.Flags().StringVarP(&configValue, "value", "v", "", "The value of the environment variable you want to set.")
	appConfigCmd.AddCommand(appConfigUnSetCmd)
	appConfigUnSetCmd.Flags().StringVarP(&envVarName, "env-var-name", "e", "", "The name of the env variable you want to unset.")
}
