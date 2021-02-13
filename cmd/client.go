package cmd

import (
	"os"

	"github.com/camdram/email-webtools/internal/client"
	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run email-webtools in client mode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("HTTP_PORT")
		token := os.Getenv("HTTP_AUTH_TOKEN")
		serverName := os.Getenv("HTTP_SERVER")
		to := os.Getenv("SMTP_TO")
		client.StartListner(port, token, serverName, to)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
