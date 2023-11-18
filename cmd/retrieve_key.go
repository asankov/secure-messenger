/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/asankov/secure-messenger/internal/secretstore"
	"github.com/spf13/cobra"
)

// exchangeKeyCmd represents the exchangeKey command
var retrieveKeyCmd = &cobra.Command{
	Use:   "retrieve-key",
	Short: "retrieve-key retrieves the key from the keychain",
	Long: `Use this command to retrieve the key from the keychain.
	
Only use this command if necessary and be careful what you do with the key after you retrieve it.`,
	Run: func(cmd *cobra.Command, args []string) {
		stdOut := cmd.OutOrStdout()
		stdErr := cmd.ErrOrStderr()

		store, err := secretstore.NewKeychainStore()
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		secretKey, err := store.Get("secret-key")
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		write(stdOut, string(secretKey))
	},
}

func init() {
	rootCmd.AddCommand(retrieveKeyCmd)
}
