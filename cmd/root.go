package cmd

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/common"
	"github.com/vumm/cli/workspace"
	"os"
)

var verbose bool
var licenseKey string
var workspacePath string

var rootCmd = &cobra.Command{
	Use:           "vumm",
	Short:         "A mod workspace for Venice Unleashed",
	Long:          "Install and manage your favourite Venice Unleashed mods.",
	SilenceErrors: true,
	SilenceUsage:  true,
	Version:       common.GetVersion(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("command failed")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initWorkspace)

	rootCmd.PersistentFlags().StringVarP(&licenseKey, "license", "l", "", "your license key")
	rootCmd.PersistentFlags().StringVarP(&workspacePath, "workspace", "w", "", "path to your server instance folder")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(publishCmd)
}

func initWorkspace() {
	log.SetHandler(cli.Default)

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if err := workspace.SetWorkspacePath(workspacePath); err != nil {
		panic(err)
	}
}
