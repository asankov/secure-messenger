/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		outErr := cmd.OutOrStderr()
		out := cmd.OutOrStdout()

		key, err := getKey()
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}
		encryptor, err := crypto.NewEncryptor(key)
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}

		msg := args[0]
		dec, err := encryptor.Decrypt(msg)
		if err != nil {
			_, _ = outErr.Write([]byte(err.Error()))
		}

		_, _ = out.Write([]byte(dec))
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
