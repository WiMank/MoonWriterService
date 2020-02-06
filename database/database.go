package database

import (
	"fmt"
)

type AppDataBaseSetting struct {
	User     string
	Password string
	Dbname   string
	Driver   string
	Sslmode  string
	ConnStr  string
}

func (appDataBaseSetting *AppDataBaseSetting) CreateConnString() {
	appDataBaseSetting.ConnStr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		appDataBaseSetting.User,
		appDataBaseSetting.Password,
		appDataBaseSetting.Dbname,
		appDataBaseSetting.Sslmode)
}
