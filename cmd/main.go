package main

import (
	"github.com/spf13/viper"

	"github.com/eyev0/timetracker/internal/config"
	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/http"
	"github.com/eyev0/timetracker/internal/log"
)

func main() {
	config.InitConfig()
	log.InitLogger(viper.GetString("LOG_LEVEL"))
	log.Logger.Infof("Starting timetracker")
	db.Init()
	http.InitHttpServer()
}
