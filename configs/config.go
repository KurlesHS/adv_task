package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBUserName  string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASS"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      int    `mapstructure:"DB_PORT"`
	DBName      string `mapstructure:"DB_NAME"`
	ServicePort int    `mapstructure:"SERVICE_PORT"`
}

func LoadConfig() (config Config, err error) {
	v := viper.New()
	v.AutomaticEnv()
	params := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_PORT", "DB_NAME", "SERVICE_PORT"}
	for _, param := range params {
		err = v.BindEnv(param)
		if err != nil {
			return
		}
	}
	err = v.Unmarshal(&config)
	if err == nil && len(config.DBName) == 0 {
		err = fmt.Errorf("error reading config")
	}
	return
}
