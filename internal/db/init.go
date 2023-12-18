package db

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/eyev0/timetracker/internal/cfg"
	"github.com/eyev0/timetracker/internal/log"
)

var Url *string = nil

func Open() (db *sqlx.DB, err error) {
	if Url == nil {
		Url = new(string)
		*Url = fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.C.PostgresUser,
			cfg.C.PostgresPassword,
			cfg.C.DatabaseHost,
			cfg.C.DatabasePort,
			cfg.C.DatabaseDb,
		)
	}
	db, err = sqlx.Open("pgx", *Url)
	if err != nil {
		log.Logger.Errorf("Unable to connect to database: %+v\n", err)
		return
	}
	return
}

func Init() {
	db, err := Open()
	if err != nil {
		log.Logger.Fatalf("Failed to open db: %+v\n", err)
	}
	defer db.Close()

	db.MustExec("SELECT 1")

	InitSchema(db)
}
