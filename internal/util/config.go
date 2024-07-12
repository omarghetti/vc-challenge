package util

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
	Environment   string `mapstructure:"ENVIRONMENT"`
	HTTPAddr      string `mapstructure:"HTTP_ADDR"`
}

func NewConfig() (Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	var config Config

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, errors.New("could not read config file")
	}

	err = viper.Unmarshal(&config)
	return config, err
}
