package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbDriver      string        `mapstructure:"DB_DRIVER"`
	DbSource      string        `mapstructure:"DB_SOURCE"`
	TestDbSource  string        `mapstructure:"DB_TEST_SOURCE"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	SymmetricKey  string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	err = viper.Unmarshal(&config)
	return
}
