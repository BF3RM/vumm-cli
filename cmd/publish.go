package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/caarlos0/ctrlc"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/pipe/archiver"
	"github.com/vumm/cli/internal/pipe/project"
	"github.com/vumm/cli/internal/pipe/publish"
	"github.com/vumm/cli/internal/pipeline"
	"time"
)

type publishCmd struct {
	cmd *cobra.Command
	tag string
}

func newPublishCmd() *publishCmd {
	root := &publishCmd{}
	root.cmd = &cobra.Command{
		Use:   "publish",
		Short: "Publish a mod",
		Long:  "Publishes a mod to the registry so that it can by installed by others.",

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("publishing...")
			start := time.Now()

			ctx, cancel := context.NewWithTimeout(30 * time.Minute)
			defer cancel()

			err := ctrlc.Default.Run(ctx, func() error {
				return pipeline.Run(ctx, project.Pipe{}, archiver.Pipe{}, publish.Pipe{Tag: root.tag})
			})
			if err != nil {
				return fmt.Errorf("failed publishing after %0.2fs: %w", time.Since(start).Seconds(), err)
			}

			log.Infof("publish success after %0.2fs", time.Since(start).Seconds())

			return nil
		},
	}

	root.cmd.Flags().StringVarP(&root.tag, "tag", "t", "latest", "version tag, latest by default")

	return root
}
