package app

import (
	"errors"

	"github.com/spf13/cobra"
)

var AppCmd = &cobra.Command{
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

// var appDomainCmd = &cobra.Command{
// 	Use:     "domain",
// 	Aliases: []string{"domains"},
// 	Short:   "Details of your application domains",
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		err := cmd.Help()
// 		if err != nil {
// 			return err
// 		}
// 		return errors.New("a valid subcommand is required")
// 	},
// }

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

func init() {
	AppCmd.AddCommand(appListCmd)
	AppCmd.AddCommand(appRemoveCmd)
	AppCmd.AddCommand(appShowCmd)
	AppCmd.AddCommand(appCreateCmd)
	appCreateCmd.Flags().StringVarP(&appSize, "size", "s", "", "Size of the application")
	appCreateCmd.Flags().StringVarP(&gitURL, "git-url", "g", "", "URL of the git repo")
	appCreateCmd.Flags().StringVarP(&image, "image", "i", "", "Container Image to pull")
	appCreateCmd.Flags().StringVarP(&branchName, "branch", "b", "", "Branch for the git repo")
	appCreateCmd.Flags().StringVarP(&tagName, "tag", "t", "", "Tag for the git repo")
	AppCmd.AddCommand(appUpdateCmd)
	appUpdateCmd.Flags().StringVarP(&appSize, "size", "s", "", "Updated Size of the application")
	appCreateCmd.Flags().StringVarP(&firewallID, "firewall-id", "f", "", "Firewall ID of the application")
	appUpdateCmd.Flags().StringVarP(&processType, "process-type", "t", "", "The type of process you want to scale. E.g. web, worker, etc.")
	appUpdateCmd.Flags().IntVarP(&processCount, "process-count", "c", 0, "The number by which you want to scale the process. E.g. 2, 3, etc.")
	// appCmd.AddCommand(appCreateCmd)
	// appCmd.AddCommand(appRemoveCmd)
	// appCmd.AddCommand(appScaleCmd)

	// //App config commands
	// appCmd.AddCommand(appConfigCmd)
	// appConfigCmd.AddCommand(appConfigShowCmd)
	// appConfigCmd.AddCommand(appConfigSetCmd)
	// appConfigSetCmd.Flags().StringVarP(&configName, "name", "n", "", "The name of the environment variable you want to set.")
	// appConfigSetCmd.Flags().StringVarP(&configValue, "value", "v", "", "The value of the environment variable you want to set.")
	// appConfigCmd.AddCommand(appConfigUnSetCmd)
	// appConfigUnSetCmd.Flags().StringVarP(&envVarName, "env-var-name", "e", "", "The name of the env variable you want to unset.")
}
