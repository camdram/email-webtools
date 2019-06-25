package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Printf("Starting Camdram Email Web Tools...")

	// Read in settings from config file.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: %s", err.Error())
	}
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mainDatabase := os.Getenv("MAIN_DB")
	serverDatabase := os.Getenv("SERVER_DB")
	token := os.Getenv("HTTP_AUTH_TOKEN")
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		log.Fatalf("Server HTTP port not set in .env file, exiting...")
	}

	// Start a webserver and listen for HTTP requests.
	driver := newSqlDriver(mysqlUser, mysqlPassword, mainDatabase, serverDatabase)
	defer driver.Clean()
	c := newController(driver, token)
	s := &http.Server{
		Addr:    ":" + port,
		Handler: c,
	}
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error starting web server: %s", err.Error())
		}
	}()
	log.Printf("Listening on port %s", port)

	// Gracefully handle SYSCALLS.
	timeout := 5 * time.Second
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down web server: %s", err.Error())
	} else {
		log.Printf("Web server terminated gracefully")
	}
}
