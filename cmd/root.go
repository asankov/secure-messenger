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
	// rootCmd.Flags().StringVar(&cfgFile, "config-file", "~/.config/secure-messenger/config.yaml", "file from which to read the config values")
	rootCmd.PersistentFlags().StringVar(&senderID, "sender-id", "", "your Sender ID")
	rootCmd.PersistentFlags().StringVar(&secretKey, "secret-key", "", "the encryption key")
	rootCmd.PersistentFlags().StringVar(&secretKeyFile, "secret-key-file", "", "the file that contains the encryption key")
}
