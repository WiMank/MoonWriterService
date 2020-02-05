package app

import "fmt"

type AlarmAppInterface interface {
	Start()
	Stop()
}

type AlarmApp struct{
	AppName string
}

func InitApp(appName string) AlarmAppInterface {
	return &AlarmApp{appName}
}

func (alarmApp *AlarmApp) Start() {
	initLogger()
	AlarmAppLog{fmt.Sprintf("====Start %s====", alarmApp.AppName), nil}.Info()
	router := initRouter()
	router.SetRoutes()
	router.ListenAndServe()
}

func (alarmApp *AlarmApp) Stop() {
	AlarmAppLog{"====Stop App====", nil}.Info()
}
