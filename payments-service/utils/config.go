package utils

import "github.com/spf13/viper"

type Config struct {
	QUEUE_NAME   string `mapstructure:"QUEUE_NAME"`
	RABBITMQ_URL string `mapstructure:"RABBITMQ_URL"`
	EXCH         string `mapstructure:"EXCH"`
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
