package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read in settings from config file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: %s", err.Error())
	}
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DB")
	connectionString := mysqlUser + ":" + mysqlPassword + "@/" + mysqlDatabase
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		log.Fatalf("Server HTTP port not set in .env file! Exiting...")
	}

	// Open a connection to the database.
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to MySQL database: %s", err.Error())
	}
	defer db.Close()

	// Prepare a query to run against the database.
	stmt, err := db.Prepare("SELECT COUNT(*) FROM queued_messages")
	if err != nil {
		log.Fatal("Error preparing SQL statement: %s", err.Error())
	}
	defer stmt.Close()

	// Serve responses using HTTP.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Postal queue length: %d", getQueueLength(stmt))
	})
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, nil)
}

func getQueueLength(stmt *sql.Stmt) int {
	var queueLength int
	err := stmt.QueryRow().Scan(&queueLength)
	if err != nil {
		log.Fatal("Error performing query: %s", err.Error())
	}
	return queueLength
}
