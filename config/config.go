package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port         string `mapstructure:"PORT"`
	GinMode      string `mapstructure:"GIN_MODE"`
	DBSource     string `mapstructure:"DATABASE_URL"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	//
	TokenSecretKey       string        `mapstructure:"TOKEN_SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(envpath string) (config ServerConfig, err error) {
	viper.AddConfigPath(envpath)
	viper.SetConfigName("local")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}
