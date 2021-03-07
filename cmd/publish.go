package cmd

import "github.com/spf13/cobra"

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a mod",
	Long:  "Publishes a mod to the registry so that it can by installed by others.",

	Run: func(cmd *cobra.Command, args []string) {

	},
}
