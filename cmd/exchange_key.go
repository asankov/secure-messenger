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
	Short: "exchange-key exchanges the secret key with the remote server",
	Long:  `This command uses a secure algorithm to exchange the secret key with the remote server over an insecure channel.`,
	Run: func(cmd *cobra.Command, args []string) {
		stdOut := cmd.OutOrStdout()
		stdErr := cmd.ErrOrStderr()

		privateKey, err := exchange.GeneratePrivateKey()
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		secretKey, err := getKey(stdErr)
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

func init() {
	rootCmd.AddCommand(exchangeKeyCmd)

	exchangeKeyCmd.Flags().StringVar(&remoteAddr, "remote-addr", "", "The address of the remote server to exchange keys with")

	_ = exchangeKeyCmd.MarkFlagRequired("remote-addr")
}
