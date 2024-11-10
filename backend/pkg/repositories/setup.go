package repositories

import "project-root/pkg/db"

var dbHandler *db.DBHandler

// SetDBHandler sets the database handler to be used by the repository
func SetDBHandler(handler *db.DBHandler) {
	dbHandler = handler
}