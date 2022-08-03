package config

import (
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port                 string        `mapstructure:"PORT"`
	GinMode              string        `mapstructure:"GIN_MODE"`
	DBSource             string        `mapstructure:"DATABASE_URL"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	TokenSecretKey       string        `mapstructure:"TOKEN_SECRET_KEY"`
	AdminSecretKey       string        `mapstructure:"ADMIN_SECRET_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(envpath string) (config ServerConfig, err error) {
	viper.AddConfigPath(envpath)
	viper.SetConfigName("remote")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
