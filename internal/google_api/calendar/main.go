package calendar

import (
	"time"

	"google.golang.org/api/calendar/v3"

	"github.com/eyev0/timetracker/internal/google_api"
	"github.com/eyev0/timetracker/internal/log"
	"github.com/eyev0/timetracker/internal/model"
)

func PostEvent(user *model.User, entry *model.Entry) (err error) {
	srv, err := google_api.GetService(user)
	if err != nil {
		log.Logger.Errorf("Unable to retrieve service: %+v", err)
		return
	}

	event := &calendar.Event{
		Summary: entry.Note,
		Start: &calendar.EventDateTime{
			DateTime: entry.StartDateTime.Format(time.RFC3339),
			TimeZone: user.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: entry.EndDateTime.Format(time.RFC3339),
			TimeZone: user.TimeZone,
		},
	}

	event, err = srv.Events.Insert(entry.CalendarId, event).Do()
	if err != nil {
		log.Logger.Errorf("Unable to insert event: %+v", err)
		return
	}
	log.Logger.Infof("Created calendar event: %+v", event)

	return
}
