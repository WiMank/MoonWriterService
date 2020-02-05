package database

import (
	"github.com/jmoiron/sqlx"
)

type AppDataBaseSetting struct {
	User     string
	Password string
	Dbname   string
	Sslmode  string
	Driver   string
}

type AppDataBaseInterface interface {
	SetUpConnection()
	OpenAppDataBase(closeDb *sqlx.DB) (*sqlx.DB, error)
	CloseAppDataBase(closeDb *sqlx.DB) bool
}
