package pkg

import "github.com/spf13/viper"

type Config struct {
	HTTP_PORT          string `mapstructure:"HTTP_PORT"`
	AUTH_GRPC_URL      string `mapstructure:"AUTH_GRPC_URL"`
	REDDIS_ADDR        string `mapstructure:"REDDIS_ADDR"`
	PAYMENT_QUEUE_NAME string `mapstructure:"PAYMENT_QUEUE_NAME"`
	RABBITMQ_URL       string `mapstructure:"RABBITMQ_URL"`
	EXCH               string `mapstructure:"EXCH"`
	POSTGRES_USER      string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD  string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_DB        string `mapstructure:"POSTGRES_DB"`
	DB_URL             string `mapstructure:"DB_URL"`
	ENCRYPTION_KEY     string `mapstructure:"ENCRYPTION_KEY"`
	MIGRATION_PATH     string `mapstructure:"MIGRATION_PATH"`
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
