package cmd

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/registry"
	"os"
)

func Execute() {
	log.SetHandler(cli.Default)
	newRootCmd().Execute()
}

type rootCmd struct {
	cmd      *cobra.Command
	registry string
	verbose  bool
}

func (cmd *rootCmd) Execute() {
	if err := cmd.cmd.Execute(); err != nil {
		log.WithError(err).Error("command failed")
		os.Exit(1)
	}
}

func newRootCmd() *rootCmd {
	root := &rootCmd{}

	root.cmd = &cobra.Command{
		Use:           "vumm",
		Short:         "A mod workspace for Venice Unleashed",
		Long:          "Install and manage your favourite Venice Unleashed mods.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       common.GetVersion(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if root.registry != "" {
				registry.SetRegistryUrl(root.registry)
			}

			if root.verbose {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	root.cmd.PersistentFlags().StringVar(&root.registry, "registry", "", "Custom registry url")
	root.cmd.PersistentFlags().BoolVarP(&root.verbose, "verbose", "v", false, "Enable verbose output")

	root.cmd.AddCommand(newInstallCmd().cmd)
	root.cmd.AddCommand(newPublishCmd().cmd)
	root.cmd.AddCommand(newUnpublishCmd().cmd)
	root.cmd.AddCommand(uninstallCmd)

	return root
}
