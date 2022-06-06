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

var appScaleCmd = &cobra.Command{
	Use:     "scale",
	Aliases: []string{"change", "modify", "upgrade"},
	Example: "civo app scale APP-NAME PROCESS-NAME=PROCESS-COUNT",
	Short:   "Scale processes of your application",
	Args:    cobra.MinimumNArgs(1),
	Run:     scaleCmd,
}

var appLogCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"log", "logs"},
	Example: "civo app log APP-NAME",
	Short:   "Check logs of your application",
	Args:    cobra.MinimumNArgs(1),
	Run:     logCmd,
}

var appSSHKeyIDCmd = &cobra.Command{
	Use:     "sshkey",
	Aliases: []string{"ssh-key-id"},
	Short:   "SSH key IDs of your application",
	Long:    "\nYou can check already existing keys with `civo sshkey ls`. If no key is found, you can create one. See `civo sshkey create --help`",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appCreateCmd)
	appCreateCmd.Flags().StringVarP(&appName, "name", "n", "", "Name of the application")
	appCreateCmd.Flags().StringVarP(&appSize, "size", "s", "", "Size of the application")
	appCreateCmd.Flags().StringVarP(&appSSHKeyIDs, "ssh-key-ids", "k", "", "SSH key IDs to authenticate to git server ")
	appCmd.AddCommand(appShowCmd)
	appCmd.AddCommand(appRemoveCmd)
	appCmd.AddCommand(appScaleCmd)
	appCmd.AddCommand(appRemoteCmd)
	appRemoteCmd.Flags().StringVarP(&remoteName, "remote-name", "r", "", "The name of remote you want to add. E.g. civo")

	//App domain commands
	appCmd.AddCommand(appDomainCmd)
	appDomainCmd.AddCommand(appDomainListCmd)
	appDomainCmd.AddCommand(appDomainAddCmd)
	appDomainCmd.AddCommand(appDomainRemoveCmd)

	//App logs
	appCmd.AddCommand(appLogCmd)

	//App config commands
	appCmd.AddCommand(appConfigCmd)
	appConfigCmd.AddCommand(appConfigShowCmd)
	appConfigCmd.AddCommand(appConfigSetCmd)
	appConfigSetCmd.Flags().StringVarP(&configName, "name", "n", "", "The name of the environment variable you want to set.")
	appConfigSetCmd.Flags().StringVarP(&configValue, "value", "v", "", "The value of the environment variable you want to set.")
	appConfigCmd.AddCommand(appConfigUnSetCmd)
	appConfigUnSetCmd.Flags().StringVarP(&envVarName, "env-var-name", "e", "", "The name of the env variable you want to unset.")

	//App SSH key commands
	appCmd.AddCommand(appSSHKeyIDCmd)
	appSSHKeyIDCmd.AddCommand(appSSHKeyAddCmd)
	appSSHKeyAddCmd.Flags().StringVarP(&appSSHKeyIDs, "ssh-key-ids", "k", "", "SSH key IDs to authenticate to git server ")
	appSSHKeyAddCmd.Flags().StringVarP(&sshKeyName, "ssh-key-name", "n", "", "The name of the SSH key you want to add.")
}
