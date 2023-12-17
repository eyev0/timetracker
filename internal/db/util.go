package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/eyev0/timetracker/internal/log"
)

func Exec(db *sqlx.DB, query string, args ...any) (result sql.Result, err error) {
	result, err = db.Exec(query, args)
	if err != nil {
		log.Logger.Errorf("Error executing query(%s) with args(%s): %+v", err)
		return
	}
	log.Logger.Debugf("exec result: %+v", result)
	return
}
