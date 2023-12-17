package db

import (
	"strings"
	"time"

	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

// modifies UpdatedAt on token
func UpsertUserToken(user *model.User, token *model.Token) (err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	token.UpdatedAt = time.Now()

	_, err = db.Exec(`INSERT INTO google_tokens
			(user_id, access_token, id_token, expires_in, refresh_token, scope, token_type)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (user_id)
			DO
				UPDATE SET access_token = $2, id_token = $3, expires_in = $4, refresh_token = $5, scope = $6, token_type = $7, updated_at = $8
		`,
		user.Id,
		token.AccessToken,
		token.IdToken,
		token.ExpiresIn,
		token.RefreshToken,
		strings.ReplaceAll(token.Scope, " ", ";"),
		token.TokenType,
		token.UpdatedAt,
	)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}

	return
}

// modifies UpdatedAt on token
func RefreshUserToken(user *model.User, token *model.Token) (err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	token.UpdatedAt = time.Now()

	_, err = db.Exec(`UPDATE google_tokens SET access_token = $1, expires_in = $2, token_type = $3, updated_at = $4 WHERE user_id = $5`,
		token.AccessToken,
		token.ExpiresIn,
		token.TokenType,
		token.UpdatedAt,
		user.Id,
	)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}

	return
}

func GetToken(user *model.User) (token *model.Token, err error) {
	db, err := Open()
	if err != nil {
		return
	}
	defer db.Close()

	row := db.QueryRowx("SELECT * FROM google_tokens WHERE user_id = $1", user.Id)

	token = new(model.Token)
	err = row.StructScan(token)
	if err != nil {
		log.Logger.Errorf("%+v", err)
		return
	}

	return
}
