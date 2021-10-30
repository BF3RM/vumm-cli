package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/internal/registry"
	"strings"
)

type grantCmd struct {
	cmd *cobra.Command
}

func newGrantCmd() *grantCmd {
	root := &grantCmd{}
	root.cmd = &cobra.Command{
		Use:   "grant <mod[@tag]> <username> <readonly|publish>",
		Short: "Grant mod permissions to a user",
		Long: `Give people mod permissions. Either grant someone with publish permissions
or give someone access to a private mod by granting them the readonly permission`,
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
			mod, tag := extractModNameAndTag(args[0])
			return registry.GrantModUserPermissions(mod, tag, args[1], args[2])
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
		Use:   "revoke <mod[@tag]> <username>",
		Short: "Revoke mod permissions of a user",
		Long:  `Revoke all mod permissions of a user`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			mod, tag := extractModNameAndTag(args[0])
			return registry.RevokeModUserPermissions(mod, tag, args[1])
		},
	}

	return root
}

func extractModNameAndTag(mod string) (string, string) {
	parts := strings.SplitN(mod, "@", 2)
	if len(parts) > 1 {
		return parts[0], parts[1]
	}

	return parts[0], ""
}
