package main

import (
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/infracstructure"
	"github.com/WiMank/MoonWriterService/interface/repository"
	"github.com/WiMank/MoonWriterService/registry"
)

func main() {
	appConfig := config.ReadConfigFile()
	infracstructure.NewLogger(appConfig)

	repository.Test()

	db := infracstructure.NewDataBase(appConfig)
	appController := registry.NewRegistry(db).NewAppController()
	infracstructure.NewRouter(appController)

}
