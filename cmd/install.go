package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/caarlos0/ctrlc"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/pipe/fetcher"
	"github.com/vumm/cli/internal/pipe/installer"
	"github.com/vumm/cli/internal/pipeline"
	"time"
)

type installCmd struct {
	cmd     *cobra.Command
	timeout time.Duration
}

func newInstallCmd() *installCmd {
	root := &installCmd{}
	root.cmd = &cobra.Command{
		Use:     "install <mod>[@<version>]",
		Aliases: []string{"add", "i"},
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

			ctx, cancel := context.NewWithTimeout(root.timeout)
			defer cancel()

			start := time.Now()
			log.Info("installing...")

			err := ctrlc.Default.Run(ctx, func() error {
				return pipeline.Run(ctx, fetcher.New(name), installer.Pipe{})
			})

			if err != nil {
				return fmt.Errorf("failed installing after %0.2fs: %w", time.Since(start).Seconds(), err)
			}

			log.Infof("successfully installed after %0.2fs", time.Since(start).Seconds())
			return nil
		},
	}

	root.cmd.Flags().DurationVar(&root.timeout, "timeout", 30*time.Minute, "Timeout for the install process")

	return root
}
