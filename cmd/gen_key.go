/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/spf13/cobra"
)

var (
	keySize int
)

// generateKeyCmd represents the generate-key command
var generateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "generate-key generates a secret key",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if keySize != 16 && keySize != 24 && keySize != 32 {
			return fmt.Errorf("key size of [%d] is not allowed. Allowed values are 16, 24 and 32.", keySize)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		stdErr := cmd.OutOrStderr()
		stdOut := cmd.OutOrStdout()

		secretKey, err := crypto.GenerateSecretKey(keySize)
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		write(stdOut, secretKey)
	},
}

func init() {
	rootCmd.AddCommand(generateKeyCmd)

	generateKeyCmd.Flags().IntVar(&keySize, "key-size", 32, "The size of the key to be generated. Allowed values are 16, 24 and 32.")
}
