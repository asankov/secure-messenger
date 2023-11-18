package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	senderID      string
	secretKey     string
	secretKeyFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "secure-messenger",
	Short: "secure-messenger allows you to send and receive encrypted messages",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&secretKey, "secret-key", "", "the encryption key. If unset the key will be read from the keychain. It is preferable to use the keychain for security reasons.")
	rootCmd.PersistentFlags().StringVar(&secretKeyFile, "secret-key-file", "", "the file that contains the encryption key. If unset the key will be read from the keychain. It is preferable to use the keychain for security reasons.")
}
