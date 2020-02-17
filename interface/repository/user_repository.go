package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/response"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type userRepository struct {
	collection *mongo.Collection
}

type UserRepository interface {
	DecodeUser(r *http.Request) domain.User
	EncodeUser(w http.ResponseWriter, response domain.UserResponse)
	InsertUser(user domain.User) response.AppResponse
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection}
}

func (ur *userRepository) DecodeUser(r *http.Request) domain.User {
	var requestUser domain.User
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Error("Decode User response! ", err)
	}
	return requestUser
}

func (ur *userRepository) EncodeUser(w http.ResponseWriter, response domain.UserResponse) {
	w.WriteHeader(response.AppResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorf("Encode User response", err)
	}
}

func (ur *userRepository) InsertUser(user domain.User) response.AppResponse {
	var localUser domain.User
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	errFind := ur.collection.FindOne(ctx, bson.D{{"user_name", user.UserName}}).Decode(&localUser)

	if errFind != nil {
		log.Info(fmt.Sprintf("Could not find user: [%s]", user.UserName))
	}

	if localUser.CheckUserExist(user) {
		return createUserExistErrorResponse(user)
	}

	_, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return createUserErrorResponse(err)
	}
	return createUserCreatedResponse(user)
}

func createUserExistErrorResponse(user domain.User) *response.UserExistError {
	userExistError := response.UserExistError{
		Message: fmt.Sprintf("User with the name [%s] is already registered", user.UserName),
		Code:    http.StatusBadRequest,
		Desc:    http.StatusText(http.StatusBadRequest),
	}
	userExistError.PrintLog(nil)
	return &userExistError
}

func createUserCreatedResponse(user domain.User) *response.UserCreated {
	userCreated := response.UserCreated{
		Message: fmt.Sprintf("User [%s] registration success!", user.UserName),
		Code:    http.StatusCreated,
		Desc:    http.StatusText(http.StatusCreated),
	}
	userCreated.PrintLog(nil)
	return &userCreated
}

func createUserErrorResponse(err error) *response.UserError {
	userError := response.UserError{
		Message: "Internal server error during user registration!",
		Code:    http.StatusInternalServerError,
		Desc:    http.StatusText(http.StatusInternalServerError),
	}
	userError.PrintLog(err)
	return &userError
}
