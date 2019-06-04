package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type SqlDriver struct {
	db   *sql.DB
	stmt *sql.Stmt
}

func newSqlDriver(mysqlUser string, mysqlPassword string, mysqlDatabase string) *SqlDriver {
	connectionString := mysqlUser + ":" + mysqlPassword + "@/" + mysqlDatabase

	// Open a connection to the database.
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Error connecting to MySQL database: %s", err.Error())
	}

	// Prepare a query to run against the database.
	stmt, err := db.Prepare("SELECT COUNT(id) as size FROM queued_messages WHERE retry_after IS NULL OR retry_after <= ADDTIME(UTC_TIMESTAMP(), '30') AND locked_at IS NULL")
	if err != nil {
		log.Fatal("Error preparing SQL statement: %s", err.Error())
	}

	return &SqlDriver{
		db:   db,
		stmt: stmt,
	}
}

func (driver *SqlDriver) GetQueueLength() int {
	var queueLength int
	if err := driver.stmt.QueryRow().Scan(&queueLength); err != nil {
		log.Fatal("Error performing query: %s", err.Error())
	}
	return queueLength
}

func (driver *SqlDriver) Clean() {
	driver.stmt.Close()
	driver.db.Close()
}
