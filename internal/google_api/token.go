package google_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

func GetOAuthToken(user *model.User) (token *oauth2.Token, err error) {
	tok, err := db.GetToken(user)
	if err != nil {
		return
	}

	token = new(oauth2.Token)
	token.AccessToken = tok.AccessToken
	token.RefreshToken = tok.RefreshToken
	token.TokenType = tok.TokenType
	token.Expiry = tok.UpdatedAt.Add(time.Duration(tok.ExpiresIn) * time.Second)

	log.Logger.Debugf("OAuth2 TOKEN: %+v", token)

	return
}

func RefreshToken(config *oauth2.Config, user *model.User, refreshToken string) (oauth2Token *oauth2.Token, err error) {
	refreshTokenURL := "https://oauth2.googleapis.com/token"
	tokenRequest := map[string]string{
		"client_id":     config.ClientID,
		"client_secret": config.ClientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	requestBody, err := json.Marshal(tokenRequest)
	if err != nil {
		log.Logger.Errorf("Error marshalling token request: %+v", err)
		return
	}

	resp, err := http.Post(refreshTokenURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Logger.Errorf("Error making the token request: %+v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		token := new(model.Token)
		if err = json.NewDecoder(resp.Body).Decode(token); err != nil {
			log.Logger.Errorf("Error decoding token response: %+v", err)
			return
		}

		err = db.RefreshUserToken(user, token)
		if err != nil {
			log.Logger.Errorf("Error upserting token to db: %+v", err)
			return
		}

		oauth2Token = new(oauth2.Token)
		oauth2Token.AccessToken = token.AccessToken
		oauth2Token.Expiry = token.UpdatedAt.Add(time.Duration(token.ExpiresIn) * time.Second)
		oauth2Token.TokenType = token.TokenType

		return
	} else {
		log.Logger.Errorf("Error occurred: %s", resp.Status)
		return nil, errors.New(fmt.Sprintf("Error: %s", resp.Status))
	}
}
