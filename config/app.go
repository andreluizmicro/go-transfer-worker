package config

import "github.com/spf13/viper"

var cfg *AppConfig

type AppConfig struct {
	NotificationClientBaseUrl    string `mapstructure:"NOTIFICATION_CLIENT_BASE_URL"`
	NotificationClientTimeout    string `mapstructure:"NOTIFICATION_CLIENT_TIMEOUT"`
	NotificationClientMaxRetries string `mapstructure:"NOTIFICATION_CLIENT_MAX_RETRIES"`
	RabbitMQConnection           string `mapstructure:"RABBITMQ_CONNECTION"`
}

func LoadConfig(path string) (*AppConfig, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}

func GetAuthorizationConfigClient() *AppConfig {
	return cfg
}
