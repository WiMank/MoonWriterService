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
