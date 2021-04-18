package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/internal/updater"
)

type updateCmd struct {
	cmd *cobra.Command
}

func newUpdateCmd() *updateCmd {
	root := &updateCmd{}
	root.cmd = &cobra.Command{
		Use:   "update",
		Short: "Check for updates",
		Long:  "Check for updates and try to self update if a newer version is available",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("checking for updates")
			updateAvailable, err := updater.CheckForUpdates()
			if err != nil {
				return err
			}

			if !updateAvailable {
				log.Info("latest version already installed")
			}

			log.Info("new version available, installing...")
			if _, err := updater.SelfUpdate(); err != nil {
				return err
			}

			return nil
		},
	}

	return root
}
