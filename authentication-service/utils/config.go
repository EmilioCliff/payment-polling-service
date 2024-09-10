package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP_PORT         string        `mapstructure:"HTTP_PORT"`
	GRPC_PORT         string        `mapstructure:"GRPC_PORT"`
	POSTGRES_USER     string        `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD string        `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_DB       string        `mapstructure:"POSTGRES_DB"`
	DB_URL            string        `mapstructure:"DB_URL"`
	HASH_COST         int           `mapstructure:"HASH_COST"`
	TOKEN_DURATION    time.Duration `mapstructure:"TOKEN_DURATION"`
	PRIVATE_KEY_PATH  string        `mapstructure:"PRIVATE_KEY_PATH"`
	PUBLIC_KEY_PATH   string        `mapstructure:"PUBLIC_KEY_PATH"`
	AUTH_QUEUE_NAME   string        `mapstructure:"AUTH_QUEUE_NAME"`
	RABBITMQ_URL      string        `mapstructure:"RABBITMQ_URL"`
	EXCH              string        `mapstructure:"EXCH"`
	ENCRYPTION_KEY    string        `mapstructure:"ENCRYPTION_KEY"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
