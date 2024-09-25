package pkg

import "github.com/spf13/viper"

type Config struct {
	SERVER_ADDRESS        string `mapstructure:"SERVER_ADDRESS"`
	RABBITMQ_URL          string `mapstructure:"RABBITMQ_URL"`
	GATEWAY_CONSUMER_NAME string `mapstructure:"GATEWAY_CONSUMER_NAME"`
	EXCH                  string `mapstructure:"EXCH"`
	EXCLUSIVE_QUEUE_NAME  string `mapstructure:"EXCLUSIVE_QUEUE_NAME"`
	AUTH_GRPC_PORT        string `mapstructure:"AUTH_GRPC_PORT"`
	AUTH_HTTP_PORT        string `mapstructure:"AUTH_HTTP_PORT"`
	PRIVATE_KEY_PATH      string `mapstructure:"PRIVATE_KEY_PATH"`
	PUBLIC_KEY_PATH       string `mapstructure:"PUBLIC_KEY_PATH"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var config Config
	err := viper.Unmarshal(&config)

	return config, err
}
