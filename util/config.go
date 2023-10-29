package util

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	PolygonApiKey string `mapstructure:"POLYGON_API_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.MergeInConfig()
	if err != nil {
		return
	}

	localConfigPath := path + "/app.env.local"
	_, err = os.Stat(localConfigPath)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if !os.IsNotExist(err) {
		viper.SetConfigName("app.env.local")
		err = viper.MergeInConfig()
		if err != nil {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
