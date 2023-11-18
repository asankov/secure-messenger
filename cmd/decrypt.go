/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt decripts an encrypted message",
	Long: `decrypt decripts an encrypted message using the secret key.

It will output the encrypted message in its JSON format on the standart output.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stdErr := cmd.OutOrStderr()
		stdOut := cmd.OutOrStdout()

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

		msg := args[0]
		dec, err := encryptor.Decrypt(msg)
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		write(stdOut, dec)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}
