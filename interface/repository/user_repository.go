package repository

import (
	"context"
	"encoding/json"
	"github.com/WiMank/AlarmService/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type userRepository struct {
	db *mongo.Client
}

type UserRepository interface {
	DecodeUser(r *http.Request) domain.User
	EncodeUser(w http.ResponseWriter, user domain.UserResponse)
	InsertUser(user domain.User) bool
	DeleteUser(user domain.User)
}

func NewUserRepository(db *mongo.Client) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) DecodeUser(r *http.Request) domain.User {
	var requestUser domain.User
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Error("Decode User error! ", err)
	}
	return requestUser
}

func (ur *userRepository) EncodeUser(w http.ResponseWriter, user domain.UserResponse) {
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Errorf("Encode User error", err)
	}
}

func (ur *userRepository) InsertUser(user domain.User) bool {
	collection := ur.db.Database("alarm_service_database").Collection("users_collection")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("InsertUser error: \n", err)
		return false
	}
	log.Info("Insert result: ", result)
	return true
}

func (ur *userRepository) DeleteUser(user domain.User) {
	collection := ur.db.Database("alarm_service_database").Collection("users_collection")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.DeleteOne(ctx, user)
	if err != nil {
		log.Errorf("DeleteUser error: \n", err)
	}
	log.Info("Delete Result: ", res)
}

func (ur *userRepository) CloseDataBase() {
	//err := ur.db.Close()
	//if err != nil {
	//	log.Errorf("Failed close database! ", err)
	//}
}
