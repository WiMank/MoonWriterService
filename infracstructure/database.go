package infracstructure

import (
	"fmt"
	"github.com/WiMank/MoonWriterService/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

//Настраиваем подключение к БД
func NewDataBase(config config.Configuration) *sqlx.DB {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s host=%s port=%d",
		config.DataBase.User,
		config.DataBase.Password,
		config.DataBase.Dbname,
		config.DataBase.Sslmode,
		config.DataBase.Host,
		config.DataBase.Port,
	)

	db, err := sqlx.Open(config.DataBase.Driver, connStr)
	if err != nil {
		panic(fmt.Errorf("Error opening database: %s \n", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("Ping error: %s \n", err))
	}

	log.Info("Successfully connected to the database!")

	return db
}
