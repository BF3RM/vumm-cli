package cmd

import "github.com/spf13/cobra"

var licenseKey string
var rootCmd = &cobra.Command{
	Use:   "vumm",
	Short: "A mod manager for Venice Unleashed",
	Long:  "Install and manage your favourite Venice Unleashed mods.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&licenseKey, "license", "l", "", "your license key")

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(publishCmd)
}
