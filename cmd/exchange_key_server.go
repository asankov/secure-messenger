/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/asankov/secure-messenger/internal/crypto/exchange"
	"github.com/spf13/cobra"
)

var (
	addr string
)

// exchangeKeyCmd represents the exchangeKey command
var exchangeKeyServerCmd = &cobra.Command{
	Use:   "exchange-key-server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stdOut := cmd.OutOrStdout()
		stdErr := cmd.ErrOrStderr()

		l, err := exchange.NewListener()
		if err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}

		stdOut.Write([]byte(fmt.Sprintf("Starting exchange key server on [%s]\n", addr)))

		if err := http.ListenAndServe(addr, l); err != nil {
			write(stdErr, err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(exchangeKeyServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exchangeKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exchangeKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	exchangeKeyServerCmd.Flags().StringVar(&addr, "addr", ":8080", "The address to listen on")
}
