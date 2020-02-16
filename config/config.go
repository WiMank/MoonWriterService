package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

//Структура конфигурации
type Configuration struct {
	DataBase struct {
		User     string
		Password string
		Dbname   string
		Driver   string
		Host     string
		Port     int
		Sslmode  string
	}
	Log struct {
		ForceColors   bool
		FullTimestamp bool
	}
}

//Читаем файл конфигурации
func ReadConfigFile() Configuration {
	var config Configuration
	viper.SetConfigName("db_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join("config"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal response config file: %s \n", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Response unmarshal config file: %s \n", err))
	}

	return config
}
