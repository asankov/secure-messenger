/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/asankov/secure-messenger/internal/crypto/exchange"
	"github.com/spf13/cobra"
)

var (
	remoteAddr string
)

// exchangeKeyCmd represents the exchangeKey command
var exchangeKeyCmd = &cobra.Command{
	Use:   "exchange-key",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stdOut := cmd.OutOrStdout()
		stdErr := cmd.ErrOrStderr()

		privateKey, err := exchange.GeneratePrivateKey()
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		secretKey, err := getKey()
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		exchanger := exchange.NewExchanger(remoteAddr)
		if err := exchanger.ExchangeSecretKey(privateKey, secretKey); err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		write(stdOut, "successfully exchanged secret key")
	},
}

func encrypt(key []byte, value []byte) ([]byte, error) {
	// implement the function

	return nil, nil
}

func init() {
	rootCmd.AddCommand(exchangeKeyCmd)

	exchangeKeyCmd.Flags().StringVar(&remoteAddr, "remote-addr", "", "The address of the remote server to exchange keys with")

	_ = exchangeKeyCmd.MarkFlagRequired("remote-addr")
}
