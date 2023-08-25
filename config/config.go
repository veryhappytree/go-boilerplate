package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env  string `mapstructure:"APP_ENV"`
	Name string `mapstructure:"APP_NAME"`
}

type DatabaseConfig struct {
	Debug    int    `mapstructure:"DB_DEBUG"`
	Host     string `mapstructure:"DB_HOST"`
	Name     string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
	Password string `mapstructure:"DB_PASSWORD"`
	User     string `mapstructure:"DB_USER"`
}

type RabbitConfig struct {
	RabbitHost string `mapstructure:"RABBIT_HOST"`
}

type RedisConfig struct {
	RedisHost string `mapstructure:"REDIS_HOST"`
}

type ServerConfig struct {
	APIDocsPort int `mapstructure:"HTTP_PORT_API_DOCS"`
	Port        int `mapstructure:"HTTP_PORT"`
}

type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Rabbit   RabbitConfig   `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
	Server   ServerConfig   `mapstructure:",squash"`
}

var LoadedConfig Config

func LoadConfig(path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic().Err(err).Msg("[APP] cannot initialize configuration")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panic().Err(err).Msg("[APP] cannot initialize configuration")
	}

	log.Info().Msg("[APP] config was loaded successfully")

	LoadedConfig = config

	return config
}
