package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/internal/registry"
)

type grantCmd struct {
	cmd *cobra.Command
}

func newGrantCmd() *grantCmd {
	root := &grantCmd{}
	root.cmd = &cobra.Command{
		Use:   "grant <mod> <user> <readonly|publish>",
		Short: "Grant mod permissions to a user",
		Long: `Give people mod permissions. Either grant someone with publish permissions
or give someone access to a private mod by granting them the readonly permission`,
		Example: "vumm grant realitymod paulhobbel readonly",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("accepts 3 arg(s), received %d", len(args))
			}

			if args[2] != "readonly" && args[2] != "publish" {
				return fmt.Errorf("invalid permission type specified, only readonly and publish are supported")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return registry.GrantModUserPermissions(args[0], args[1], args[2])
		},
	}

	return root
}

type revokeCmd struct {
	cmd *cobra.Command
}

func newRevokeCmd() *revokeCmd {
	root := &revokeCmd{}
	root.cmd = &cobra.Command{
		Use:     "revoke <mod> <user>",
		Short:   "Revoke mod permissions of a user",
		Long:    `Revoke all mod permissions of a user`,
		Example: "vumm grant realitymod paulhobbel readonly",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return registry.RevokeModUserPermissions(args[0], args[1])
		},
	}

	return root
}
