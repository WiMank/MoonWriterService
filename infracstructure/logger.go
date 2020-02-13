package infracstructure

import (
	"github.com/WiMank/AlarmService/config"
	log "github.com/sirupsen/logrus"
)

//Настраиваем цвет вывода и временные метки логгера
func NewLogger(configuration config.Configuration) {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   configuration.Log.ForceColors,
		FullTimestamp: configuration.Log.FullTimestamp,
	})
}
