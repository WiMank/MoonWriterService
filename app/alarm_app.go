package app

import "fmt"

type AlarmAppInterface interface {
	Start()
	Stop()
}

type AlarmApp struct {
	AppName string
	Router  RouterInterface
}

func InitApp(appName string) AlarmAppInterface {
	initLogger()
	return &AlarmApp{appName, initRouter()}
}

func (alarmApp *AlarmApp) Start() {
	AlarmAppLog{fmt.Sprintf("====Start %s====", alarmApp.AppName), nil}.Info()
	alarmApp.Router.SetRoutes()
	alarmApp.Router.ListenAndServe()
}

func (alarmApp *AlarmApp) Stop() {
	AlarmAppLog{"====Stop App====", nil}.Info()
}
