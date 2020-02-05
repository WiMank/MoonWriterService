package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
		AlarmAppLog{Message: "Cannot open database", Err: err}.Fatal()
	}
	appDataBaseSetting.db = db
}

func (appDataBaseSetting *AppDataBaseSetting) CloseAppDataBase() {
	closeErr := appDataBaseSetting.db.Close()
	if closeErr != nil {
		AlarmAppLog{Message: "Cannot close database", Err: closeErr}.Fatal()
	}
}
