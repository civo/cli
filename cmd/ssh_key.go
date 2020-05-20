package cmd

import (
	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:     "sshkey",
	Aliases: []string{"ssh", "ssh-key", "sshkeys"},
	Short:   "Details of Civo SSH keys",
}

func init() {
	rootCmd.AddCommand(sshKeyCmd)
	sshKeyCmd.AddCommand(sshKeyListCmd)
	sshKeyCmd.AddCommand(sshKeyCreateCmd)
	sshKeyCmd.AddCommand(sshKeyRemoveCmd)

	sshKeyCreateCmd.Flags().StringVarP(&keyCreate, "key", "k", "", "The path of the key")
	sshKeyCreateCmd.MarkFlagRequired("key")
}
