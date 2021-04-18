package cmd

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vumm/cli/internal/common"
	"github.com/vumm/cli/internal/registry"
	"github.com/vumm/cli/internal/updater"
	"os"
)

func Execute() {
	log.SetHandler(cli.Default)
	newRootCmd().Execute()
}

type rootCmd struct {
	cmd     *cobra.Command
	verbose bool
}

func (cmd *rootCmd) Execute() {
	if err := cmd.cmd.Execute(); err != nil {
		log.WithError(err).Error("command failed")
		os.Exit(1)
	}
}

func newRootCmd() *rootCmd {
	cobra.OnInitialize(initConfig)

	root := &rootCmd{}

	root.cmd = &cobra.Command{
		Use:           "vumm",
		Short:         "A mod workspace for Venice Unleashed",
		Long:          "Install and manage your favourite Venice Unleashed mods.",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       common.GetFullVersion(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			registryUrl := viper.GetString("registry")
			if registryUrl != "" {
				registry.SetRegistryUrl(registryUrl)
			}

			token := viper.GetString("token")
			if token != "" {
				registry.SetRegistryAccessToken(token)
			}

			if root.verbose {
				log.SetLevel(log.DebugLevel)
			}
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			release, available, _ := updater.PeriodicCheckForUpdates()
			if available {
				log.WithField("version", release.Version()).Info("update available, run vumm update to install")
			}
		},
	}

	root.cmd.PersistentFlags().String("registry", "", "Custom registry url")
	root.cmd.PersistentFlags().String("token", "", "A access token to access the registry")
	root.cmd.PersistentFlags().BoolVarP(&root.verbose, "verbose", "v", false, "Enable verbose output")

	cobra.CheckErr(viper.BindPFlag("registry", root.cmd.PersistentFlags().Lookup("registry")))
	cobra.CheckErr(viper.BindPFlag("token", root.cmd.PersistentFlags().Lookup("token")))

	root.cmd.AddCommand(newInstallCmd().cmd)
	//root.cmd.AddCommand(uninstallCmd)

	root.cmd.AddCommand(newPublishCmd().cmd)
	root.cmd.AddCommand(newUnpublishCmd().cmd)
	root.cmd.AddCommand(newPackCmd().cmd)

	root.cmd.AddCommand(newLoginCmd().cmd)
	root.cmd.AddCommand(newRegisterCmd().cmd)
	root.cmd.AddCommand(newGrantCmd().cmd)
	root.cmd.AddCommand(newRevokeCmd().cmd)

	root.cmd.AddCommand(newUpdateCmd().cmd)

	return root
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("json")
	viper.SetConfigName(".vumm")
	viper.SetEnvPrefix("vumm")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cobra.CheckErr(err)
		}
	}
}
