package app

import (
	"../database"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type AlarmApp struct {
	AppName            string
	Router             AlarmAppRouter
	AppDataBaseSetting database.AppDataBaseSetting
}

func InitApp(appName string) AlarmApp {
	initLogger()

	return AlarmApp{appName, initRouter(), setUpDbConnection()}
}

func setUpDbConnection() database.AppDataBaseSetting {
	log.Info("<<<Set Up DataBase Connection>>>")
	connection := database.AppDataBaseSetting{
		User:     "postgres",
		Password: "korabl",
		Dbname:   "alarm_database",
		Driver:   "postgres",
		Sslmode:  "disable",
	}
	connection.CreateConnString()
	return connection
}

func (alarmApp *AlarmApp) Start() {
	log.Info(fmt.Sprintf("====Start %s====", alarmApp.AppName))
	alarmApp.Router.SetRoutes(alarmApp.AppDataBaseSetting)
	alarmApp.Router.ListenAndServe()
}

func (alarmApp *AlarmApp) Stop() {
	log.Info(fmt.Sprintf("====Stop %s====", alarmApp.AppName))
}
