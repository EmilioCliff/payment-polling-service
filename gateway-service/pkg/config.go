package pkg

import "github.com/spf13/viper"

type Config struct {
	SERVER_ADDRESS       string `mapstructure:"SERVER_ADDRESS"`
	RABBITMQ_URL         string `mapstructure:"RABBITMQ_URL"`
	QUEUE_NAME           string `mapstructure:"QUEUE_NAME"`
	EXCH                 string `mapstructure:"EXCH"`
	EXCLUSIVE_QUEUE_NAME string `mapstructure:"EXCLUSIVE_QUEUE_NAME"`
	AUTH_GRPC_PORT string `mapstructure:"AUTH_GRPC_PORT"`
	AUTH_HTTP_PORT string `mapstructure:"AUTH_HTTP_PORT"`
	PRIVATE_KEY_PATH  string        `mapstructure:"PRIVATE_KEY_PATH"`
	PUBLIC_KEY_PATH   string        `mapstructure:"PUBLIC_KEY_PATH"`
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
