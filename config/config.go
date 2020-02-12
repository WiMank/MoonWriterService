package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

type Configuration struct {
	DataBase struct {
		User     string
		Password string
		Dbname   string
		Driver   string
		Addr     string
		Sslmode  bool
	}
}

func ReadConfigFile() Configuration {
	var config Configuration
	viper.SetConfigName("db_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join("$GOPATH", "src", "github.com", "WiMank", "AlarmService", "config"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		panic(fmt.Errorf("Error unmarshal config file: %s \n", err))
	}
	return config
}
