package server

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type SQLDriver struct {
	mainDb    *sql.DB
	serverDb  *sql.DB
	queueStmt *sql.Stmt
	heldStmt  *sql.Stmt
}

func newSQLDriver(mysqlUser string, mysqlPassword string, mainDatabase string, serverDatabase string) (*SQLDriver, error) {
	mainConnString := mysqlUser + ":" + mysqlPassword + "@/" + mainDatabase
	serverConnString := mysqlUser + ":" + mysqlPassword + "@/" + serverDatabase

	// Open a connection to the database.
	mainDb, err := sql.Open("mysql", mainConnString)
	if err != nil {
		return nil, err
	}
	serverDb, err := sql.Open("mysql", serverConnString)
	if err != nil {
		return nil, err
	}

	// Prepare queries to run against the database.
	queueStmt, err := mainDb.Prepare("SELECT COUNT(id) as size FROM queued_messages WHERE retry_after IS NULL OR retry_after <= ADDTIME(UTC_TIMESTAMP(), '30') AND locked_at IS NULL")
	if err != nil {
		return nil, err
	}
	heldStmt, err := serverDb.Prepare("SELECT COUNT(id) as size FROM messages WHERE held = 1")
	if err != nil {
		return nil, err
	}

	return &SQLDriver{
		mainDb:    mainDb,
		serverDb:  serverDb,
		queueStmt: queueStmt,
		heldStmt:  heldStmt,
	}, nil
}

func (driver *SQLDriver) GetQueueLength() (int, error) {
	var queueLength int
	if err := driver.queueStmt.QueryRow().Scan(&queueLength); err != nil {
		return 0, err
	}
	return queueLength, nil
}

func (driver *SQLDriver) GetHeldMessageCount() (int, error) {
	var heldMessageCount int
	if err := driver.heldStmt.QueryRow().Scan(&heldMessageCount); err != nil {
		return 0, err
	}
	return heldMessageCount, nil
}

func (driver *SQLDriver) Clean() {
	driver.queueStmt.Close()
	driver.heldStmt.Close()
	driver.mainDb.Close()
	driver.serverDb.Close()
}
