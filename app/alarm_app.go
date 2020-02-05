package app

import (
	"../database"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type AlarmAppInterface interface {
	Start()
	Stop()
}

type AlarmApp struct {
	AppName            string
	Router             RouterInterface
	AppDataBaseSetting database.AppDataBaseInterface
}

func InitApp(appName string) AlarmAppInterface {
	initLogger()
	return &AlarmApp{appName, initRouter(), setUpDbConnection()}
}

func setUpDbConnection() database.AppDataBaseInterface {
	log.Info("<<<Set Up DataBase Connection>>>")
	return &database.AppDataBaseSetting{
		User:     "postgres",
		Password: "korabl",
		Dbname:   "alarm_database",
		Driver:   "postgres",
		Sslmode:  false,
	}
}

func (alarmApp *AlarmApp) Start() {
	log.Info(fmt.Sprintf("====Start %s====", alarmApp.AppName))
	alarmApp.Router.SetRoutes()
	alarmApp.Router.ListenAndServe()
}

func (alarmApp *AlarmApp) Stop() {
	log.Info(fmt.Sprintf("====Stop %s====", alarmApp.AppName))
}
