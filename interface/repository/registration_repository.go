package repository

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type registrationRepository struct {
	db              *sqlx.DB
	responseCreator response.AppResponseCreator
	validator       *validator.Validate
}

type RegistrationRepository interface {
	DecodeRequest(r *http.Request) request.UserRegistrationRequest
	InsertUser(request request.UserRegistrationRequest) response.AppResponse
}

func NewUserRepository(
	db *sqlx.DB,
	responseCreator response.AppResponseCreator,
	validator *validator.Validate,
) RegistrationRepository {
	return &registrationRepository{db, responseCreator, validator}
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
	var exist bool
	existQuery := "SELECT EXISTS (SELECT FROM users WHERE user_name=$1 AND user_pass=$2)::boolean"
	err := ur.db.QueryRowx(existQuery, authReq.User.UserName, authReq.User.UserPass).Scan(&exist)
	if err != nil {
		log.Error("findUserEntity: ", err)
		return true
	}
	return exist
}

func (ur *registrationRepository) insertUserEntity(entity *domain.UserEntity) bool {
	insertQuery := "INSERT INTO users (user_name, user_pass, user_role, premium) VALUES ($1, $2, 'user', false)"
	result, err := ur.db.Exec(insertQuery, entity.UserName, entity.UserPass)
	if err != nil {
		return false
	}

	if result != nil {
		return true
	}

	return false
}
