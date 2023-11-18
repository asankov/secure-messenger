/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/asankov/secure-messenger/internal/secretstore"
	"github.com/spf13/cobra"
)

var (
	keySize        int
	outputToStdout bool
)

// generateKeyCmd represents the generate-key command
var generateKeyCmd = &cobra.Command{
	Use:   "generate-key",
	Short: "generate-key generates a secret key",
	Long: `generate-key generates a secret key, which is used to encrypt and decrypt messages.

By default, the command will save the key into the OS keychain.
This is the best option for security reasons.

If instead, you want to output it to the terminal (or save it to file), you can use the "--output-to-stdout" flag and redirect the output to a file.

The key size can be specified with the --key-size flag. Allowed values are 16, 24 and 32.
	`,
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

		if outputToStdout {
			write(stdErr, "Consider using the keychain to store the key, so that it don't get lost or exposed.")
			write(stdOut, secretKey)

			return
		}

		store, err := secretstore.NewKeychainStore()
		if err != nil {
			// write the errors to stderr, but the key to stdout so that it can be used in a pipe
			write(stdErr, "error while creating keychain store: "+err.Error())
			write(stdErr, "Outputting the key to stdout instead, so that it does not get lost")
			write(stdOut, secretKey)
			write(stdErr, err.Error())
			os.Exit(1)
		}

		if _, err := store.StoreSecretKey(secretKey); err != nil {
			// write the errors to stderr, but the key to stdout so that it can be used in a pipe
			write(stdErr, "error while storing key into the keychain: "+err.Error())
			write(stdErr, "Outputting the key to stdout instead, so that it does not get lost")
			write(stdOut, secretKey)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateKeyCmd)

	generateKeyCmd.Flags().IntVar(&keySize, "key-size", 32, "The size of the key to be generated. Allowed values are 16, 24 and 32.")
	generateKeyCmd.Flags().BoolVar(&outputToStdout, "output-to-stdout", false, "Output the key to stdout. It is prefered to keep this to false, and use the keychain for security reasons.")
}
