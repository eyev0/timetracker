package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/eyev0/timetracker/internal/cfg"
	"github.com/eyev0/timetracker/internal/db"
	"github.com/eyev0/timetracker/internal/http"
	"github.com/eyev0/timetracker/internal/log"
)

func main() {
	err := cfg.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error loading config file: %w", err))
	}
	log.InitLogger(viper.GetString("LOG_LEVEL"))
	log.Logger.Infof("Starting timetracker")
	db.Init()
	http.InitServer()
}
