package repository

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/WiMank/MoonWriterService/interface/utils"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type registrationRepository struct {
	collectionUsers *mongo.Collection
	responseCreator response.AppResponseCreator
	validator       *validator.Validate
}

type RegistrationRepository interface {
	DecodeRequest(r *http.Request) request.UserRegistrationRequest
	InsertUser(request request.UserRegistrationRequest) response.AppResponse
}

func NewUserRepository(
	collectionUsers *mongo.Collection,
	responseCreator response.AppResponseCreator,
	validator *validator.Validate,
) RegistrationRepository {
	return &registrationRepository{collectionUsers, responseCreator, validator}
}

func (ur *registrationRepository) DecodeRequest(r *http.Request) request.UserRegistrationRequest {
	var requestUser request.UserRegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Error("Decode User response error! ", err)
	}

	return requestUser
}

func (ur *registrationRepository) InsertUser(request request.UserRegistrationRequest) response.AppResponse {
	if request.ValidateRequest(ur.validator) {

		userExist := ur.findUserEntity(request)

		if !userExist {
			insertComplete := ur.insertUserEntity(&request.User)
			if insertComplete {
				return ur.responseCreator.CreateResponse(response.UserCreatedResponse{}, request.User.UserName)
			} else {
				return ur.responseCreator.CreateResponse(response.UserInsertErrorResponse{}, request.User.UserName)
			}
		} else {
			return ur.responseCreator.CreateResponse(response.UserExistResponse{}, request.User.UserName)
		}
	}

	return ur.responseCreator.CreateResponse(response.ValidateErrorResponse{}, "")
}

func (ur *registrationRepository) findUserEntity(authReq request.UserRegistrationRequest) bool {
	count, err := ur.collectionUsers.CountDocuments(utils.GetContext(),
		bson.D{
			{"user_name", authReq.User.UserName},
			{"user_pass", authReq.User.UserPass},
		})

	if err != nil {
		return false
	}

	if count != 1 {
		return false
	}

	return true
}

func (ur *registrationRepository) insertUserEntity(entity *domain.UserEntity) bool {
	_, err := ur.collectionUsers.InsertOne(utils.GetContext(),
		bson.D{
			{"user_name", entity.UserName},
			{"user_pass", entity.UserPass},
			{"user_role", "user"},
			{"is_premium_user", false},
		})

	if err != nil {
		return false
	}

	return true
}
