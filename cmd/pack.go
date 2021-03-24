package cmd

import (
	"fmt"
	"github.com/apex/log"
	"github.com/caarlos0/ctrlc"
	"github.com/spf13/cobra"
	"github.com/vumm/cli/internal/context"
	"github.com/vumm/cli/internal/pipe/archiver"
	"github.com/vumm/cli/internal/pipe/project"
	"github.com/vumm/cli/internal/pipeline"
	"time"
)

type packCmd struct {
	cmd     *cobra.Command
	timeout time.Duration
}

func newPackCmd() *packCmd {
	root := &packCmd{}
	root.cmd = &cobra.Command{
		Use:   "pack",
		Short: "Create a tarball for a mod",
		Long:  "Packs the current mod into a tarball like the publish command would do",

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("packing...")
			start := time.Now()

			ctx, cancel := context.NewWithTimeout(root.timeout)
			defer cancel()

			err := ctrlc.Default.Run(ctx, func() error {
				return pipeline.Run(ctx, project.Pipe{}, archiver.Pipe{Store: true})
			})
			if err != nil {
				return fmt.Errorf("failed packing after %0.2fs: %w", time.Since(start).Seconds(), err)
			}

			log.Infof("pack success after %0.2fs", time.Since(start).Seconds())

			return nil
		},
	}

	root.cmd.Flags().DurationVar(&root.timeout, "timeout", 30*time.Minute, "Timeout for the pack process")

	return root
}
