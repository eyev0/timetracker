package google_api

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

func GetService(user *model.User) (srv *calendar.Service, err error) {
	ctx := context.Background()
	creds, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Logger.Errorf("Unable to read client secret file: %+v", err)
		return
	}

	config, err := google.ConfigFromJSON(creds, calendar.CalendarEventsScope)
	if err != nil {
		log.Logger.Errorf("Unable to parse client secret file to config: %+v", err)
		return
	}

	oauth2Token, err := GetOAuthToken(user)
	if err != nil {
		log.Logger.Errorf("Unable to retrieve user token: %+v", err)
		return
	}

	if !oauth2Token.Valid() {
		log.Logger.Infof("Token expired, refreshing")

		oauth2Token, err = RefreshToken(config, user, oauth2Token.RefreshToken)
		if err != nil {
			log.Logger.Errorf("Unable to update token: %+v", err)
			return
		}

		log.Logger.Infof("New token: %+v", oauth2Token)

		// token, err = config.TokenSource(ctx, token).Token()
		// if err != nil {
		// 	log.Logger.Errorf("Unable to update token: %+v", err)
		// 	return
		// }
	}

	srv, err = calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, oauth2Token)))
	if err != nil {
		log.Logger.Errorf("Unable to create new service: %+v", err)
		return
	}

	return
}
