package app

import (
	"fmt"
	"github.com/WiMank/AlarmService/database"
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
	AlarmAppLog{"<<<Set Up DataBase Connection>>>", nil}.Info()
	return &database.AppDataBaseSetting{
		User:     "postgres",
		Password: "korabl",
		Dbname:   "alarm_database",
		Driver:   "postgres",
		Sslmode:  false,
	}
}

func (alarmApp *AlarmApp) Start() {
	AlarmAppLog{fmt.Sprintf("====Start %s====", alarmApp.AppName), nil}.Info()
	alarmApp.Router.SetRoutes()
	alarmApp.Router.ListenAndServe()
}

func (alarmApp *AlarmApp) Stop() {
	AlarmAppLog{"====Stop App====", nil}.Info()
}
