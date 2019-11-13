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
	log.Println("Starting Email Web Tools in server mode")
	driver, err := newSQLDriver(mysqlUser, mysqlPassword, mainDatabase, serverDatabase)
	if err != nil {
		log.Fatalln("Failed to initialise connection to database:", err)
	}
	defer driver.Clean()
	c := newController(driver, token)
	s := &http.Server{
		Addr:    ":" + port,
		Handler: c,
	}
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("Failed to start web server:", err)
		}
	}()
	log.Println("Listening on port", port)

	// Gracefully handle SYSCALLS.
	timeout := 5 * time.Second
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln("Failed to terminate web server:", err)
	} else {
		log.Println("Web server terminated gracefully")
		log.Println("Bye-bye!")
	}
}
