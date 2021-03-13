package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vumm/cli/publish"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a mod",
	Long:  "Publishes a mod to the registry so that it can by installed by others.",

	Run: func(cmd *cobra.Command, args []string) {
		publisher, err := publish.NewPublisher()
		if err != nil {
			cobra.CheckErr(err)
		}
		publisher.Publish()

		//ignorer, err := publish.CompileFileIgnorerCwd()
		//if err != nil {
		//	cobra.CheckErr(err)
		//}
		//
		//// TEST
		//cwd, err := os.Getwd()
		//if err != nil {
		//	cobra.CheckErr(err)
		//}
		//
		//format := publish.NewTarGZPackager()
		//format.SetIgnorer(ignorer)
		//err = format.Make(cwd, "archive.tar.gz")
		//if err != nil {
		//	cobra.CheckErr(err)
		//}
	},
}
