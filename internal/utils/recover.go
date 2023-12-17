package utils

import "github.com/eyev0/timetracker/internal/log"

func Recover_gracefully(message string) {
	if r := recover(); r != nil {
		log.Logger.Warnf("Recovered with message - %s: %+v", message, r)
	}
}
