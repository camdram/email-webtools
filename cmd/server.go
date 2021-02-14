package cmd

import (
	"github.com/camdram/email-webtools/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run email-webtools in server mode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()
		port := viper.GetString("HTTP_PORT")
		token := viper.GetString("HTTP_AUTH_TOKEN")
		mysqlUser := viper.GetString("MYSQL_USER")
		mysqlPassword := viper.GetString("MYSQL_PASSWORD")
		mainDatabase := viper.GetString("MAIN_DB")
		serverDatabase := viper.GetString("SERVER_DB")
		server.StartServer(port, token, mysqlUser, mysqlPassword, mainDatabase, serverDatabase)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
