package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type AppDataBaseSetting struct {
	User     string
	Password string
	Dbname   string
	Driver   string
	Sslmode  bool
	connStr  string
	db       *sqlx.DB
}

type AppDataBaseInterface interface {
	SetUpConnection()
	OpenAppDataBase()
	CloseAppDataBase()
}

func (appDataBaseSetting *AppDataBaseSetting) SetUpConnection() {
	appDataBaseSetting.connStr = fmt.Sprintf("%s %s %s %t",
		appDataBaseSetting.User,
		appDataBaseSetting.Password,
		appDataBaseSetting.Dbname, appDataBaseSetting.Sslmode)
}

func (appDataBaseSetting *AppDataBaseSetting) OpenAppDataBase() {
	db, err := sqlx.Open(appDataBaseSetting.Driver, appDataBaseSetting.connStr)
	if err != nil {
		log.Fatal("Cannot open database", err)
		return
	}
	appDataBaseSetting.db = db
}

func (appDataBaseSetting *AppDataBaseSetting) CloseAppDataBase() {
	err := appDataBaseSetting.db.Close()
	if err != nil {
		log.Fatal("Cannot close database", err)
	}
}
