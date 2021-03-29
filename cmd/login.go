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
	cmd *cobra.Command
}

func newLoginCmd() *loginCmd {
	root := &loginCmd{}
	root.cmd = &cobra.Command{
		Use: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			token, err := registry.Login(strings.TrimSpace(username), strings.TrimSpace(string(bytePassword)), registry.AccessTokenTypeReadonly)
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

	return root
}
