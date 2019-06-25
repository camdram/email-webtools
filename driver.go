package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type SqlDriver struct {
	mainDb    *sql.DB
	serverDb  *sql.DB
	queueStmt *sql.Stmt
	heldStmt  *sql.Stmt
}

func newSqlDriver(mysqlUser string, mysqlPassword string, mainDatabase string, serverDatabase string) *SqlDriver {
	mainConnString := mysqlUser + ":" + mysqlPassword + "@/" + mainDatabase
	serverConnString := mysqlUser + ":" + mysqlPassword + "@/" + serverDatabase

	// Open a connection to the database.
	mainDb, err := sql.Open("mysql", mainConnString)
	serverDb, err := sql.Open("mysql", serverConnString)
	if err != nil {
		log.Fatalf("Error connecting to MySQL database: %s", err.Error())
	}

	// Prepare queries to run against the database.
	queueStmt, err := mainDb.Prepare("SELECT COUNT(id) as size FROM queued_messages WHERE retry_after IS NULL OR retry_after <= ADDTIME(UTC_TIMESTAMP(), '30') AND locked_at IS NULL")
	heldStmt, err := serverDb.Prepare("SELECT COUNT(id) as size FROM messages WHERE held = 1")
	if err != nil {
		log.Fatal("Error preparing SQL statement: %s", err.Error())
	}

	return &SqlDriver{
		mainDb:    mainDb,
		serverDb:  serverDb,
		queueStmt: queueStmt,
		heldStmt:  heldStmt,
	}
}

func (driver *SqlDriver) GetQueueLength() int {
	var queueLength int
	if err := driver.queueStmt.QueryRow().Scan(&queueLength); err != nil {
		log.Fatal("Error performing query: %s", err.Error())
	}
	return queueLength
}

func (driver *SqlDriver) GetHeldMessageCount() int {
	var heldMessageCount int
	if err := driver.heldStmt.QueryRow().Scan(&heldMessageCount); err != nil {
		log.Fatal("Error performing query: %s", err.Error())
	}
	return heldMessageCount
}

func (driver *SqlDriver) Clean() {
	driver.queueStmt.Close()
	driver.heldStmt.Close()
	driver.mainDb.Close()
	driver.serverDb.Close()
}
