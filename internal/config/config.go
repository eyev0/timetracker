package config

import (
	"fmt"
	"time"

	config "github.com/spf13/viper"
)

func InitConfig() {
	config.SetConfigName("app")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}
}

type Config struct {
	FrontEndOrigin string `mapstructure:"FRONTEND_ORIGIN"`

	JWTTokenSecret string        `mapstructure:"JWT_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRES_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`

	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectUrl string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
}

func LoadConfig(path string) (c Config, err error) {
	config.AddConfigPath(path)
	config.SetConfigType("env")
	config.SetConfigName("auth")

	config.AutomaticEnv()

	err = config.ReadInConfig()
	if err != nil {
		return
	}

	err = config.Unmarshal(&c)
	return
}
