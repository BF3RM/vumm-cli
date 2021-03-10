package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vumm/cli/workspace"
)

var licenseKey string
var workspacePath string

var rootCmd = &cobra.Command{
	Use:   "vumm",
	Short: "A mod workspace for Venice Unleashed",
	Long:  "Install and manage your favourite Venice Unleashed mods.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initWorkspace)

	rootCmd.PersistentFlags().StringVarP(&licenseKey, "license", "l", "", "your license key")
	rootCmd.PersistentFlags().StringVarP(&workspacePath, "workspace", "w", "", "path to your server instance folder")

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(publishCmd)
}

func initWorkspace() {
	if err := workspace.SetWorkspacePath(workspacePath); err != nil {
		panic(err)
	}
}
