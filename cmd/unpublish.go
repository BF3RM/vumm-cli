package cmd

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

type unpublishCmd struct {
	cmd *cobra.Command
}

func newUnpublishCmd() *unpublishCmd {
	root := &unpublishCmd{}
	root.cmd = &cobra.Command{
		Use:   "unpublish <mod>@<version>",
		Short: "Remove a mod from the registry",
		Long:  "Removes a mod version from the registry.",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			parts := strings.Split(args[0], "@")
			if len(parts) != 2 {
				return fmt.Errorf("invalid mod version specified")
			}
			mod := parts[0]
			version, err := semver.NewVersion(parts[1])
			if err != nil {
				return fmt.Errorf("invalid mod version specified")
			}

			start := time.Now()
			log.Info("unpublishing...")

			if _, err = client.Mods.UnpublishModVersion(cmd.Context(), mod, version); err != nil {
				return fmt.Errorf("unpublish unsuccessful after %0.2fs: %v", time.Since(start).Seconds(), err)
			}

			log.Infof("unpublish successful after %0.2fs", time.Since(start).Seconds())
			return nil
		},
	}

	return root
}
