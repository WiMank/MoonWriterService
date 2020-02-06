package controller

import (
	"../database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type BaseController struct {
	DbSetting database.AppDataBaseSetting
}

type BaseControllerInterface interface {
	OpenAppDataBase()
	CloseAppDataBase()
}

func (c *BaseController) OpenAppDataBase() *sqlx.DB {
	db, err := sqlx.Open(c.DbSetting.Driver, c.DbSetting.ConnStr)
	if err != nil {
		log.Fatal("Cannot open database ", err)
		return nil
	}
	return db
}

func (c *BaseController) CloseAppDataBase(db *sqlx.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal("Cannot close database ", err)
	}
}

func (c *BaseController) CloseRows(row *sqlx.Rows) {
	err := row.Close()
	if err != nil {
		log.Fatal("Cannot close row ", err)
	}
}
