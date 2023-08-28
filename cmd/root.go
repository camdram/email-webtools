package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logFile string
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "email-webtools",
	Short: "A tiny micro-service to ensure that Camdram can send & receive emails 24/7/365",
	Long:  "",
}

func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to config file (default is ENV)")
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "l", "", "path to log file (default is standard output)")
}

func initConfig() {
	if os.Getenv("JOURNAL_STREAM") != "" {
		log.Default().SetFlags(0)
	}
}

func loadConfig() {
	viper.AutomaticEnv()
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln("Unable to read config file:", err)
		}
	}
	ensureConfig()
}

func ensureConfig() {
	if viper.GetString("HTTP_PORT") == "" {
		log.Fatalln("Server HTTP port not set")
	}
	if viper.GetString("HTTP_AUTH_TOKEN") == "" {
		log.Fatalln("Server HTTP auth token not set")
	}
}
