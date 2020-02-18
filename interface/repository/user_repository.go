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
	collection      *mongo.Collection
	responseCreator response.AppResponseCreator
}

type UserRepository interface {
	DecodeRequest(r *http.Request) request.UserRegistrationRequest
	InsertUser(request request.UserRegistrationRequest) response.AppResponse
}

func NewUserRepository(collection *mongo.Collection, responseCreator response.AppResponseCreator) UserRepository {
	return &userRepository{collection, responseCreator}
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

	findUserErr := ur.collection.FindOne(ctx, bson.D{{"user_name", request.User.UserName}}).Decode(&localUserEntity)
	if findUserErr != nil {
		log.Errorf(fmt.Sprintf("User [%s] not found: %v", request.User.UserName, findUserErr))
	}

	if localUserEntity.CheckUserExist(request.User) {
		return ur.responseCreator.CreateResponse(response.UserExistResponse{}, request.User.UserName)
	}

	if _, err := ur.collection.InsertOne(ctx, request.User); err != nil {
		return ur.responseCreator.CreateResponse(response.UserInsertErrorResponse{}, request.User.UserName)
	}

	return ur.responseCreator.CreateResponse(response.UserCreatedResponse{}, request.User.UserName)
}
