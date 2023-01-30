package cmd

import (
	"errors"

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

func init() {
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appCreateCmd)
	appCreateCmd.Flags().StringVarP(&appName, "name", "n", "", "Name of the application")
	appCreateCmd.Flags().StringVarP(&appSize, "size", "s", "", "Size of the application")
	appCmd.AddCommand(appRemoveCmd)
	appCmd.AddCommand(appScaleCmd)
	appScaleCmd.Flags().StringVarP(&processType, "process-type", "t", "", "The type of process you want to scale. E.g. web, worker, etc.")
	appScaleCmd.Flags().IntVarP(&processCount, "process-count", "c", 0, "The number by which you want to scale the process. E.g. 2, 3, etc.")

	//App config commands
	appCmd.AddCommand(appConfigCmd)
	appConfigCmd.AddCommand(appConfigShowCmd)
	appConfigCmd.AddCommand(appConfigSetCmd)
	appConfigSetCmd.Flags().StringVarP(&configName, "name", "n", "", "The name of the environment variable you want to set.")
	appConfigSetCmd.Flags().StringVarP(&configValue, "value", "v", "", "The value of the environment variable you want to set.")
	appConfigCmd.AddCommand(appConfigUnSetCmd)
	appConfigUnSetCmd.Flags().StringVarP(&envVarName, "env-var-name", "e", "", "The name of the env variable you want to unset.")
}
