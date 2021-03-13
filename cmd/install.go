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

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		i, err := installer.NewInstaller()
		if err != nil {
			cobra.CheckErr(err)
		}

		installedMods := i.GetInstalledMods()
		fmt.Printf("Currently there are %d mod(s) installed\n", len(installedMods))
		for _, installedMod := range installedMods {
			fmt.Printf("\t%s - %s\n", installedMod.Name, installedMod.Version)
		}
		fmt.Println()

		if i.HasMissingMods() {
			missingMods := i.GetMissingMods()
			fmt.Printf("Detected %d missing mod(s)\n", len(missingMods))

			for _, missingMod := range missingMods {
				fmt.Printf("\t%s - %s\n", missingMod.Name, missingMod.VersionConstraints)
			}
			fmt.Println()
		}

		if err = i.InstallMod(name); err != nil {
			cobra.CheckErr(err)
		}
	},
}
