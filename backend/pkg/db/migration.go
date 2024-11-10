package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DBHandler holds the instances for main and message databases
type DBHandler struct {
	MainDB *sql.DB
	MsgDB  *sql.DB
}

// NewDBHandler creates a new DBHandler instance
func NewDBHandler() *DBHandler {
	return &DBHandler{}
}

// InitMainDB initializes the main database
func (handler *DBHandler) InitMainDB(dataSourceName, initScript string) error {
	var err error
	handler.MainDB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	if err = handler.MainDB.Ping(); err != nil {
		return err
	}

	return runInitScript(handler.MainDB, initScript)
}

// InitMsgDB initializes the messages database
func (handler *DBHandler) InitMsgDB(dataSourceName string) error {
	var err error
	handler.MsgDB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	return handler.MsgDB.Ping()
}

// runInitScript executes the initialization script for the database
func runInitScript(db *sql.DB, scriptPath string) error {
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(script))
	return err
}
