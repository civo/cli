package cmd

import (
	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:     "sshkey",
	Aliases: []string{"ssh", "ssh-key"},
	Short:   "Details of Civo Ssh Keys",
}

func init() {
	rootCmd.AddCommand(sshKeyCmd)
	sshKeyCmd.AddCommand(sshKeyListCmd)
	sshKeyCmd.AddCommand(sshKeyCreateCmd)
	sshKeyCmd.AddCommand(sshKeyRemoveCmd)

	/*
		Flags for ssh key create
	*/
	sshKeyCreateCmd.Flags().StringVarP(&keyCreate, "key", "k", "", "The path of the key")
	sshKeyCreateCmd.MarkFlagRequired("key")
}
