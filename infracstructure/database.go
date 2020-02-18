package infracstructure

import (
	"context"
	"fmt"
	"github.com/WiMank/MoonWriterService/config"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

//Настраиваем подключение к БД
func NewDataBase(config config.Configuration) *mongo.Database {
	connStr := fmt.Sprintf(
		"mongodb://%s:%d",
		config.DataBase.Host,
		config.DataBase.Port,
	)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		panic(fmt.Errorf("Connect to database response: %s \n", err))
	}

	errPing := client.Ping(ctx, readpref.Primary())
	if errPing != nil {
		panic(fmt.Errorf("Ping response: %s \n", err))
	}

	log.Info("Successfully connected to the database!")

	return client.Database(config.DataBase.Dbname)
}
