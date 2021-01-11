package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/camdram/email-webtools/internal/client"
	"github.com/camdram/email-webtools/internal/server"
	"github.com/joho/godotenv"
)

var version = "dev-edge"
var copyright = "Copyright (c) 2019–2021 The Association of Cambridge Theatre Societies and contributing groups."

var port, token, mysqlUser, mysqlPassword, mainDatabase, serverDatabase, serverName, to string

func main() {
	// Parse command line flags.
	verThenExit := flag.Bool("version", false, "print software version and exit")
	logFile := flag.String("log", "", "path to log file")
	confFile := flag.String("config", ".env", "path to config file")
	mode := flag.String("mode", "", "Program mode")
	flag.Parse()

	// Print version and exit if the user asked
	if *verThenExit {
		path := os.Args[0]
		fmt.Println(path, "version", version)
		fmt.Println(copyright)
		os.Exit(0)
	}

	// Deal with logging output.
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	} else {
		log.SetFlags(0)
	}

	// Read in environmental variable configuration.
	readConfig(confFile)

	// Run the actual application.
	if *mode == "server" {
		server.StartServer(port, token, mysqlUser, mysqlPassword, mainDatabase, serverDatabase)
	} else if *mode == "client" {
		client.StartListner(port, token, serverName, to)
	} else {
		fmt.Println("Need to specify '--mode server' or '--mode client'")
		os.Exit(1)
	}
}

func readConfig(confFile *string) {
	if err := godotenv.Load(*confFile); err != nil {
		log.Fatalln(err)
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
	to = os.Getenv("SMTP_TO")
}
