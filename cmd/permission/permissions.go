package permission

import (
	"errors"

	"github.com/spf13/cobra"
)

// PermissionsCmd manages permissions
var PermissionsCmd = &cobra.Command{
	Use:     "permissions",
	Aliases: []string{"permission"},
	Short:   "List available permissions",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	PermissionsCmd.AddCommand(permissionsListCmd)
}
