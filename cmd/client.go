package cmd

import (
	"github.com/camdram/email-webtools/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run email-webtools in client mode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()
		userAgent := "email-webtools/" + version + " (+https://github.com/camdram/email-webtools)"
		port := viper.GetString("HTTP_PORT")
		token := viper.GetString("HTTP_AUTH_TOKEN")
		serverName := viper.GetString("HTTP_SERVER")
		to := viper.GetString("SMTP_TO")
		client.StartListner(port, token, serverName, userAgent, to)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
