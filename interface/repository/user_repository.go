package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
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
	DecodeRequest(r *http.Request) request.UserRegistrationRequest
	InsertUser(request request.UserRegistrationRequest) response.AppResponse
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection}
}

func (ur *userRepository) DecodeRequest(r *http.Request) request.UserRegistrationRequest {
	var requestUser request.UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Error("Decode User response error! ", err)
	}
	return requestUser
}

func (ur *userRepository) InsertUser(request request.UserRegistrationRequest) response.AppResponse {
	var localUserEntity domain.UserEntity
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	errFind := ur.collection.FindOne(ctx, bson.D{{"user_name", request.User.UserName}}).Decode(&localUserEntity)
	if errFind != nil {
		log.Info(fmt.Sprintf("Could not find user: [%s]", request.User.UserName))
	}

	if localUserEntity.CheckUserExist(request.User) {
		return createUserExistErrorResponse(request.User)
	}

	if _, err := ur.collection.InsertOne(ctx, request.User); err != nil {
		return createUserErrorResponse(err)
	}
	return createUserCreatedResponse(request.User)
}

func createUserExistErrorResponse(user domain.UserEntity) *response.UserExistResponse {
	userExistError := response.UserExistResponse{
		Message: fmt.Sprintf("User with the name [%s] is already registered", user.UserName),
		Code:    http.StatusBadRequest,
		Desc:    http.StatusText(http.StatusBadRequest),
	}
	userExistError.PrintLog(nil)
	return &userExistError
}

func createUserCreatedResponse(user domain.UserEntity) *response.UserCreatedResponse {
	userCreated := response.UserCreatedResponse{
		Message: fmt.Sprintf("User [%s] registration success!", user.UserName),
		Code:    http.StatusCreated,
		Desc:    http.StatusText(http.StatusCreated),
	}
	userCreated.PrintLog(nil)
	return &userCreated
}

func createUserErrorResponse(err error) *response.UserInsertErrorResponse {
	userError := response.UserInsertErrorResponse{
		Message: "Internal server error during user registration!",
		Code:    http.StatusInternalServerError,
		Desc:    http.StatusText(http.StatusInternalServerError),
	}
	userError.PrintLog(err)
	return &userError
}
