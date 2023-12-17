package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/utils"
)

var schema = `
create table users
(
    id          uuid primary key      default gen_random_uuid(),
    name        varchar(100),
    email       varchar(100) unique,
    password    varchar(100) not null,
    time_zone   varchar(50)  not null default 'Europe/Moscow',
    calendar_id varchar(200)          default 'primary',
    created_at  timestamptz           default now(),
    updated_at  timestamptz           default now()
);

create table entries
(
    id              uuid primary key default gen_random_uuid(),
    user_id         uuid references users (id),
    start_timestamp timestamptz      default now(),
    end_timestamp   timestamptz,
    note            varchar(100),
    calendar_id     varchar(200)
);

create table google_tokens
(
    id            uuid primary key default gen_random_uuid(),
    user_id       uuid unique references users (id),
    access_token  varchar(2000),
    id_token      varchar(2000),
    expires_in    int,
    refresh_token varchar(200),
    scope         varchar(2000),
    token_type    varchar(100),
    created_at    timestamptz      default now(),
    updated_at    timestamptz      default now()
);
`

func InitSchema(db *sqlx.DB) {
	defer utils.Recover_gracefully("Could not create schema")
	db.MustExec(schema)
	log.Logger.Infof("Successfully created database schema")
}
