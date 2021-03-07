package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall [mod]",
	Aliases: []string{"remove"},
	Short:   "Remove a installed mod",
	Long:    "This command uninstalls a mods and any leftover dependencies",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("install requires at least a name of the mod")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

	},
}
