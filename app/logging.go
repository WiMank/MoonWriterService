package app

import log "github.com/sirupsen/logrus"

type LoggerInterface interface {
	Info()
	Warn()
	Error()
	Fatal()
}

type AlarmAppLog struct {
	Message string
	Err     error
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.Info("<<<Init Logger>>>")
}

func (appLog AlarmAppLog) Info() {
	log.Info(appLog.Message)
}

func (appLog AlarmAppLog) Warn() {
	log.Warn(appLog.Message, " error: ", appLog.Err)
}

func (appLog AlarmAppLog) Error() {
	log.Error(appLog.Message, " error: ", appLog.Err)
}

func (appLog AlarmAppLog) Fatal() {
	log.Fatal(appLog.Message, " error: ", appLog.Err)
}
