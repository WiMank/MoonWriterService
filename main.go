package main

import (
	"github.com/WiMank/AlarmService/config"
	"github.com/WiMank/AlarmService/infracstructure"
	"github.com/WiMank/AlarmService/registry"
)

func main() {
	appConfig := config.ReadConfigFile()
	infracstructure.NewLogger(appConfig)
	db := infracstructure.NewDataBase(appConfig)
	appController := registry.NewRegistry(db).NewAppController()
	infracstructure.NewRouter(appController)
}
