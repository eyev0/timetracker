package db

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

// modifies input entry with values selected from db
func CreateEntry(entry *model.Entry, user *model.User) (err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	row := db.QueryRowx(
		"INSERT INTO entries (id, user_id, note, calendar_id) VALUES ($1, $2, $3, $4) RETURNING *",
		entry.Id,
		entry.UserId,
		entry.Note,
		user.CalendarId,
	)
	row.StructScan(entry)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}
	log.Logger.Debugf("CreateEntry result: %+v", entry)

	return
}

// modifies input, so entry should not be nil
func GetEntryById(tx *sqlx.Tx, entry *model.Entry) (err error) {
	row := tx.QueryRowx("SELECT * FROM entries WHERE id = $1", entry.Id)
	err = row.StructScan(entry)
	log.Logger.Debugf("GetEntryById result: %+v", entry)

	return
}

func GetCurrentUserEntry(user *model.User) (entry *model.Entry, err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	row := db.QueryRowx("SELECT * FROM entries WHERE user_id = $1 AND end_timestamp is null", user.Id)
	entry = new(model.Entry)
	err = row.StructScan(entry)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Logger.Infof("%s", err)
		} else {
			log.Logger.Errorf("%v", err)
		}
	}
	log.Logger.Debugf("GetCurrentUserEntry result: %+v", entry)

	return
}

func UpdateEntry(user *model.User, input *model.UpdateEntryInput) (entry *model.Entry, tx *sqlx.Tx, err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	tx = db.MustBegin()

	template := "UPDATE entries SET end_timestamp = %s, note = %s WHERE user_id = $1 AND %s RETURNING *"
	var end_timestamp, note, whereClause string
	var varargs []any = make([]any, 0, 3)

	varargs = append(varargs, user.Id)

	if input.EndDateTime != nil {
		end_timestamp = fmt.Sprintf("$%d", len(varargs)+1)
		varargs = append(varargs, input.EndDateTime)
	} else {
		end_timestamp = "now()"
	}

	if input.Note != nil {
		note = fmt.Sprintf("$%d", len(varargs)+1)
		varargs = append(varargs, input.Note)
	} else {
		note = "note"
	}

	if input.Id != nil {
		whereClause = fmt.Sprintf("id = $%d", len(varargs)+1)
		varargs = append(varargs, input.Id)
	} else {
		whereClause = "end_timestamp IS NULL"
	}

	stmt := fmt.Sprintf(template, end_timestamp, note, whereClause)
	row := tx.QueryRowx(
		stmt,
		varargs...,
	)

	entry = new(model.Entry)
	err = row.StructScan(entry)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Logger.Infof("%s", err)
		} else {
			log.Logger.Errorf("%v", err)
		}
	}

	return
}
