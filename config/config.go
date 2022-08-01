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
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			config = ServerConfig{
				Port:                 viper.GetString("PORT"),
				GinMode:              viper.GetString("GIN_MODE"),
				DBSource:             viper.GetString("DATABASE_URL"),
				MigrationURL:         viper.GetString("MIGRATION_URL"),
				TokenSecretKey:       viper.GetString("TOKEN_SECRET_KEY"),
				AccessTokenDuration:  viper.GetDuration("ACCESS_TOKEN_DURATION"),
				RefreshTokenDuration: viper.GetDuration("REFRESH_TOKEN_DURATION"),
			}
			return
		}
		return
	}
	err = viper.Unmarshal(&config)
	return
}
