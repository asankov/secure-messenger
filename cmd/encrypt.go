/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	Short: "secure-messenger allows you to send and received encrypted messages",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		outErr := cmd.OutOrStderr()
		out := cmd.OutOrStdout()

		msg, err := messages.NewMessage(senderID, receiverID, payload)
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}

		key, err := getKey()
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}
		encryptor, err := crypto.NewEncryptor(key)
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}

		json, err := msg.ToJSON()
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}

		enc, err := encryptor.Encrypt(json)
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}

		_, _ = out.Write([]byte(enc))
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	encryptCmd.Flags().StringVar(&receiverID, "receiver-id", "", "the ID of the person receiving the message")
	encryptCmd.Flags().StringVar(&payload, "payload", "", "the payload to be encrypted")
}
