package utils

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	PostgresDSN string `mapstructure:"POSTGRES_DSN"`
}

var config AppConfig

func GetConfig() *AppConfig {
	return &config
}

func init() {
	vp := viper.New()
	vp.AutomaticEnv()
	vp.SetDefault("POSTGRES_DSN", "host=localhost user=postgres password=12345 dbname=postgres port=15432 sslmode=disable")
	if err := vp.Unmarshal(&config); err != nil {
		Logger.Info("Not able to unmarshall config.", zap.Error(err))

	}
}
