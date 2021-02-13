package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	logFilePath  string
	confFilePath string
)

var rootCmd = &cobra.Command{
	Use:   "email-webtools",
	Short: "A tiny micro-service to ensure that Camdram can send & receive emails 24/7/365",
	Long:  "",
}

func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	loadEnvironment()
	ensureEnvironment()
}

func initConfig() {
	rootCmd.Flags().StringVarP(&logFilePath, "log", "l", "", "path to log file")
	rootCmd.Flags().StringVarP(&confFilePath, "environment", "e", "/etc/email-webtools/.env", "path to environment file")
}

func loadEnvironment() {
	if err := godotenv.Load(confFilePath); err != nil {
		log.Fatalf("Failed to load environment file: %v", err)
	}
}

func ensureEnvironment() {
	if os.Getenv("HTTP_PORT") == "" {
		log.Fatalln("Server HTTP port not set in .env file, exiting...")
	}
	if os.Getenv("HTTP_AUTH_TOKEN") == "" {
		log.Fatalln("Server HTTP auth token not set in .env file, exiting...")
	}
}

func configureLogging() {
	if logFilePath != "" {
		f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		// Log to standard output.
		log.SetFlags(0)
	}
}
