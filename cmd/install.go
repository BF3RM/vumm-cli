package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/workspace"
)

var installCmd = &cobra.Command{
	Use:     "install [mod] [version]",
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

		graph := workspace.NewModDependencyGraph(workspace.ResolveModDependencyFromString(name))
		resolved, errs := graph.Resolve()
		if !resolved {
			cobra.CheckErr(errs)
		}

		resolvedMods := graph.GetResolvedDependencies()
		fmt.Printf("Resolved %d dependencies\n", len(resolvedMods))
		for _, resolvedMod := range resolvedMods {
			fmt.Printf("\t%s - %s\n", resolvedMod.Name, resolvedMod.Version)
		}

		//modVersion, err := registry.GetModVersion(name, version)
		//if err != nil {
		//	cobra.CheckErr(err)
		//}
		//fmt.Println(modVersion)

		//print(workspace.GetInstalledMods())
	},
}
