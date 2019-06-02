package main

import (
	"fmt"
	"os"
	"log"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DB")
	connectionString := mysqlUser + ":" + mysqlPassword + "@/" + mysqlDatabase

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Error connecting to MySQL database: ", err.Error())
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT COUNT(*) FROM queued_messages")
	if err != nil {
		log.Fatal("Error preparing SQL statement: ", err.Error())
	}
	defer stmtOut.Close()

	var queueLength int
	err = stmtOut.QueryRow().Scan(&queueLength)
	if err != nil {
		log.Fatal("Error performing query: ", err.Error())
	}

	fmt.Printf("The length of the Postal mail queue is: %d \n", queueLength)
}
