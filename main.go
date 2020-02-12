package main

import (
	"github.com/WiMank/AlarmService/config"
	log "github.com/sirupsen/logrus"
)

func main() {

	appConfig := config.ReadConfigFile()

	log.Info(appConfig.DataBase)

}
