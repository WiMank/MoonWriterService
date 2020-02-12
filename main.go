package main

import (
	"github.com/WiMank/AlarmService/config"
	"github.com/WiMank/AlarmService/infracstructure"
)

func main() {
	appConfig := config.ReadConfigFile()
	infracstructure.NewLogger(appConfig)
	infracstructure.NewDataBase(appConfig)
}
