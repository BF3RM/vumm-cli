package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/installer"
)

var installCmd = &cobra.Command{
	Use:     "install [mod]",
	Aliases: []string{"add"},
	Short:   "Install a mod",
	Long:    "This command installs a mod and any mods that it depends on.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("install requires at least a name of the mod")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		i, err := installer.NewInstaller()
		if err != nil {
			return err
		}

		if err = i.InstallMod(name); err != nil {
			return err
		}

		return nil
	},
}
