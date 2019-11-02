package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(port string, token string, mysqlUser string, mysqlPassword string, mainDatabase string, serverDatabase string) {
	log.Printf("Starting Camdram Email Web Tools server...")
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
