package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	DBSource     string `mapstructure:"DATABASE_URL"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	//
	TokenSecretKey       string        `mapstructure:"TOKEN_SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadServerConfig() (config ServerConfig, err error) {
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("MIGRATION_URL")
	viper.BindEnv("TOKEN_SECRET_KEY")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("REFRESH_TOKEN_DURATION")
	err = viper.Unmarshal(&config)
	return
}
