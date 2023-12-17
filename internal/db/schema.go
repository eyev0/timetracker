package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/eyev0/timetracker/internal/utils"
)

var schema = `
create table users (
    id uuid primary key,
    name varchar(100),
    email varchar(100) unique,
    password varchar(100) not null,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

create table entries (
    id BIGSERIAL primary key,
    user_id BIGSERIAL references users(id),
    start_timestamp timestamptz default now(),
    end_timestamp timestamptz,
    note varchar(100),
    calendar_id varchar(200)
);
`

func InitSchema(db *sqlx.DB) {
	defer utils.Recover_gracefully("Could not create schema")
	db.MustExec(schema)
}
