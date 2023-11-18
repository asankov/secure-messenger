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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
