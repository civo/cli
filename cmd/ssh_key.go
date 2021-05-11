package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:     "sshkey",
	Aliases: []string{"ssh", "ssh-key", "sshkeys"},
	Short:   "Details of Civo SSH keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	rootCmd.AddCommand(sshKeyCmd)
	sshKeyCmd.AddCommand(sshKeyListCmd)
	sshKeyCmd.AddCommand(sshKeyCreateCmd)
	sshKeyCmd.AddCommand(sshKeyRemoveCmd)

	sshKeyCreateCmd.Flags().StringVarP(&keyCreate, "key", "k", "", "The path of the key")
	sshKeyCreateCmd.MarkFlagRequired("key")
}
