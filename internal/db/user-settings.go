package db

import (
	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

func UpdateUserSettings(user *model.User, settings *model.SettingsInput) (err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	result, err := db.Exec("UPDATE users SET calendar_id = $1", settings.CalendarId)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}
	log.Logger.Debugf("UpdateUserSettings result: %+v", result)

	return
}
