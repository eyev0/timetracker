package db

import (
	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

func CreateUser(user *model.User) (err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}
	log.Logger.Debugf("Insert result: %+v", res)

	return nil
}

func GetUserByEmail(email string) (user *model.User, err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	row  := db.QueryRowx("SELECT * from users WHERE email = $1", email)
	user = new(model.User)
	err = row.StructScan(user)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}
	log.Logger.Debugf("GetUserByEmail result: %+v", user)

	return
}

func GetUserById(id string) (user *model.User, err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	row  := db.QueryRowx("SELECT * from users WHERE id = $1", id)
	user = new(model.User)
	err = row.StructScan(user)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}
	log.Logger.Debugf("GetUserById result: %+v", user)

	return
}
