package cmd

// sshkey list -- list all SSH keys [ls, all]
// sshkey upload NAME FILENAME -- upload the SSH public key in FILENAME to a new key called NAME [create, new]
// sshkey remove ID -- remove the SSH public key with ID [delete, destroy, rm]

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
}
