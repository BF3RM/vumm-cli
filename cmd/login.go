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

type loginCmd struct {
	tokenType string
	cmd       *cobra.Command
}

func newLoginCmd() *loginCmd {
	root := &loginCmd{}
	root.cmd = &cobra.Command{
		Use:   "login",
		Short: "Login with username and password",
		Long:  "Login with a username and password so you can publish your mod",
		RunE: func(cmd *cobra.Command, args []string) error {
			tokenType := registry.AccessTokenTypePublish
			if root.tokenType == "readonly" {
				tokenType = registry.AccessTokenTypeReadonly
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

			fmt.Println("Logging in...")
			token, err := registry.Login(strings.TrimSpace(username), strings.TrimSpace(string(bytePassword)), tokenType)
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
