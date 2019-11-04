package main

import (
	"log"
	"os"

	"github.com/camdram/email-webtools/internal/client"
	"github.com/camdram/email-webtools/internal/server"
	"github.com/joho/godotenv"
)

var port, token, mysqlUser, mysqlPassword, mainDatabase, serverDatabase, serverName string

func readConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: %s", err.Error())
	}
	port = os.Getenv("HTTP_PORT")
	if port == "" {
		log.Fatalf("Server HTTP port not set in .env file, exiting...")
	}
	token = os.Getenv("HTTP_AUTH_TOKEN")
	if token == "" {
		log.Fatalf("Server HTTP auth token not set in .env file, exiting...")
	}
	mysqlUser = os.Getenv("MYSQL_USER")
	mysqlPassword = os.Getenv("MYSQL_PASSWORD")
	mainDatabase = os.Getenv("MAIN_DB")
	serverDatabase = os.Getenv("SERVER_DB")
	serverName = os.Getenv("HTTP_SERVER")
}

func main() {
	logfile, err := os.OpenFile("ewt.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	readConfig()
	if len(os.Args) > 1 && os.Args[1] == "--server" {
		server.StartServer(port, token, mysqlUser, mysqlPassword, mainDatabase, serverDatabase)
	} else {
		client.StartListner(port, token, serverName)
	}
}
