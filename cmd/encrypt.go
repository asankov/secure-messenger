/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/asankov/secure-messenger/internal/messages"
	"github.com/spf13/cobra"
)

var (
	receiverID string
	payload    string
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt generate an encrypted message",
	Long: `encrypt generate a JSON message that has sender id, receiver id, payload and a timestamp and encrypts the message with the secret key.

"sender-id", "receiver-id" and "payload" are required flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		stdOut := cmd.OutOrStdout()
		stdErr := cmd.OutOrStderr()

		msg, err := messages.NewMessage(senderID, receiverID, payload)
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		key, err := getKey(stdErr)
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}
		encryptor, err := crypto.NewEncryptor(key)
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		json, err := msg.ToJSON()
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		enc, err := encryptor.Encrypt(json)
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		write(stdOut, enc)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVar(&senderID, "sender-id", "", "your Sender ID")
	encryptCmd.Flags().StringVar(&receiverID, "receiver-id", "", "the ID of the person receiving the message")
	encryptCmd.Flags().StringVar(&payload, "payload", "", "the payload to be encrypted")

	_ = encryptCmd.MarkFlagRequired("sender-id")
	_ = encryptCmd.MarkFlagRequired("receiver-id")
	_ = encryptCmd.MarkFlagRequired("payload")
}
