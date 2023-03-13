package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbHost     string `mapstructure:"DB_HOST"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbUrl      string `mapstructure:"DB_URL"`
	DbName     string `mapstructure:"DB_NAME"`
}

func LoadConfig() (config *Config, err error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// enables system environment variables take precedence over the ones
	// from environment variable files.
	viper.AutomaticEnv()

	// tells viper to start reading the config
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	return

}
