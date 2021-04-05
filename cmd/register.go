package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vumm/cli/internal/registry"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

type registerCmd struct {
	tokenType string
	cmd       *cobra.Command
}

func newRegisterCmd() *registerCmd {
	root := &registerCmd{}
	root.cmd = &cobra.Command{
		Use:   "register",
		Short: "Register with username and password",
		Long:  "Register with a username and password so you can publish your mod or access private mods",
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenType, err := registry.ParseTokenType(root.tokenType)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Username: ")
			username, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			fmt.Print("Password: ")
			bytePassword, err := terminal.ReadPassword(syscall.Stdin)
			if err != nil {
				return err
			}
			fmt.Println()

			fmt.Println("Registering...")
			token, err := registry.Register(strings.TrimSpace(username), strings.TrimSpace(string(bytePassword)), tokenType)
			if err != nil {
				return err
			}

			viper.Set("token", token.Token)
			if err = viper.WriteConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					return viper.SafeWriteConfig()
				} else {
					return err
				}
			}

			return nil
		},
	}

	root.cmd.Flags().StringVar(&root.tokenType, "type", "publish", "Type of token to be generated (publish, readonly)")

	return root
}