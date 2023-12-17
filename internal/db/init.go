package db

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	config "github.com/spf13/viper"

	"github.com/eyev0/timetracker/internal/log"
)

var Url *string = nil

func Open() (db *sqlx.DB, err error) {
	if Url == nil {
		Url = new(string)
		*Url = fmt.Sprintf(
			"postgres://%s:%s@%s/%s",
			config.GetString("POSTGRES_USER"),
			config.GetString("POSTGRES_PASSWORD"),
			config.GetString("DATABASE_ADDRESS"),
			config.GetString("DATABASE_DB"),
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

	// InitSchema(db)
}
