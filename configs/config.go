package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBUserName string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASS"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig() (config Config, err error) {
	v := viper.New()
	v.AutomaticEnv()
	v.BindEnv("DB_USER")
	v.BindEnv("DB_PASS")
	v.BindEnv("DB_HOST")
	v.BindEnv("DB_PORT")
	v.BindEnv("DB_NAME")
	err = v.Unmarshal(&config)
	if err == nil && len(config.DBName) == 0 {
		err = fmt.Errorf("error reading config")
	}
	return
}
