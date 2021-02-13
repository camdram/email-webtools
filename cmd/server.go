package cmd

import (
	"os"

	"github.com/camdram/email-webtools/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run email-webtools in server mode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("HTTP_PORT")
		token := os.Getenv("HTTP_AUTH_TOKEN")
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		mainDatabase := os.Getenv("MAIN_DB")
		serverDatabase := os.Getenv("SERVER_DB")
		server.StartServer(port, token, mysqlUser, mysqlPassword, mainDatabase, serverDatabase)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
