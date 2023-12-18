package cfg

import (
	"time"

	"github.com/spf13/viper"
)

func LoadConfig() (err error) {
	v := viper.New()

	v.SetConfigFile("./.env")

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	v.AutomaticEnv()

	err = v.Unmarshal(&C)
	return
}

var C *Config

type Config struct {
	ServerIP   string `mapstructure:"SERVER_IP"`
	ServerPort int    `mapstructure:"SERVER_PORT"`

	LogFilepath string `mapstructure:"LOG_FILEPATH"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`

	DatabaseHost string `mapstructure:"DATABASE_HOST"`
	DatabasePort int    `mapstructure:"DATABASE_PORT"`
	DatabaseDb   string `mapstructure:"DATABASE_DB"`

	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`

	FrontEndOrigin string `mapstructure:"FRONTEND_ORIGIN"`

	JWTTokenSecret string        `mapstructure:"JWT_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRES_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`

	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectUrl string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
}
