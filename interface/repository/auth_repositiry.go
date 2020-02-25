package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/WiMank/MoonWriterService/interface/utils"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type authRepository struct {
	collectionUsers    *mongo.Collection
	collectionSessions *mongo.Collection
	responseCreator    response.AppResponseCreator
}

type AuthRepository interface {
	DecodeRequest(r *http.Request) request.AuthenticateUserRequest
	AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse
}

func NewAuthRepository(collectionUsers *mongo.Collection, collectionSessions *mongo.Collection, responseCreator response.AppResponseCreator) AuthRepository {
	return &authRepository{collectionUsers, collectionSessions, responseCreator}
}

func (ar *authRepository) DecodeRequest(r *http.Request) request.AuthenticateUserRequest {
	var requestUser request.AuthenticateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Errorf("Decode User response error:\n", err)
	}
	return requestUser
}

func (ar *authRepository) AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse {
	localUserEntity, errFindUser := ar.findUserEntity(authReq)
	if errFindUser != nil {
		return ar.responseCreator.CreateResponse(response.UserFindResponse{}, authReq.User.UserName)
	}

	session := ar.findSession(authReq.MobileKey)

	if localUserEntity.CheckUserNameAndPass(authReq.User) {
		ar.checkSessionsCount(authReq)

		access, errAuthenticate := createAccessToken(localUserEntity)
		refresh, errAuthenticate := createRefreshToken(localUserEntity)

		if errAuthenticate != nil {
			return ar.responseCreator.CreateResponse(response.TokenErrorResponse{}, authReq.User.UserName)
		}

		if (session != nil) && (session.CheckMkExist(authReq.MobileKey)) {
			updateResult, updateErr := ar.updateSession(access, refresh, localUserEntity, authReq)
			if updateErr != nil {
				ar.responseCreator.CreateResponse(response.SessionUpdateFailedResponse{}, authReq.User.UserName)
			}
			return ar.responseCreator.CreateResponse(response.TokenResponse{
				Message:      fmt.Sprintf("Tokens updated for [%s]", localUserEntity.UserName),
				SessionId:    updateResult,
				AccessToken:  access,
				RefreshToken: refresh}, authReq.User.UserName)

		} else {
			insertResult, insertErr := ar.insertSession(access, refresh, authReq, localUserEntity)
			if insertErr != nil {
				ar.responseCreator.CreateResponse(response.SessionInsertFailedResponse{}, authReq.User.UserName)
			}
			return ar.responseCreator.CreateResponse(response.TokenResponse{
				Message:      fmt.Sprintf("Tokens created for [%s]", localUserEntity.UserName),
				SessionId:    insertResult,
				AccessToken:  access,
				RefreshToken: refresh}, authReq.User.UserName)
		}
	}
	return ar.responseCreator.CreateResponse(response.UnauthorizedResponse{}, authReq.User.UserName)
}

func (ar *authRepository) findUserEntity(authReq request.AuthenticateUserRequest) (*domain.UserEntity, error) {
	var localUserEntity domain.UserEntity
	userBson := bson.M{"user_name": authReq.User.UserName}
	if errFind := ar.collectionUsers.FindOne(utils.GetContext(), userBson).Decode(&localUserEntity); errFind != nil {
		return nil, errFind
	}
	return &localUserEntity, nil
}

func (ar *authRepository) findSession(mk string) *domain.SessionEntity {
	var session domain.SessionEntity
	errMk := ar.collectionSessions.FindOne(utils.GetContext(), bson.M{"mobile_key": mk}).Decode(&session)
	if errMk != nil {
		log.Error("FindSession error: ", errMk)
		return nil
	}
	return &session
}

func (ar *authRepository) checkSessionsCount(authReq request.AuthenticateUserRequest) {
	userBson := bson.M{"user_name": authReq.User.UserName}
	count, errCount := ar.collectionSessions.CountDocuments(utils.GetContext(), userBson)
	if count > 5 {
		ar.clearSessions(userBson)
	}
	if errCount != nil {
		log.Errorf("CheckSessionsCount error: \n", errCount)
	}
}

func (ar *authRepository) clearSessions(userBson bson.M) {
	result, errDelete := ar.collectionSessions.DeleteMany(utils.GetContext(), userBson)
	if errDelete != nil {
		log.Errorf("ClearSessions error:\n", errDelete)
	}
	if result != nil {
		log.Info("Delete result: ", result.DeletedCount)
	}
}

func createAccessToken(entity *domain.UserEntity) (string, error) {
	tokenTime := utils.GetAccessTokenTime()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "Moon Writer",
		"user": entity.Id,
		"role": entity.UserRole,
		"exp":  tokenTime,
	})
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Errorf("Access token signed error:\n", err)
		return "", err
	}
	return tokenString, nil
}

func createRefreshToken(entity *domain.UserEntity) (string, error) {
	tokenTime := utils.GetRefreshTokenTime()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "Moon Writer",
		"user": entity.Id,
		"exp":  tokenTime,
	})
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Errorf("Refresh token signed error:\n", err)
		return "", err
	}
	return tokenString, nil
}

func (ar *authRepository) insertSession(access string, refresh string, authReq request.AuthenticateUserRequest, entity *domain.UserEntity) (string, error) {
	newSession := createSession(access, refresh, authReq, entity)
	insertResult, errInsert := ar.collectionSessions.InsertOne(utils.GetContext(), bson.D{
		{"user_id", newSession.UserId},
		{"user_name", newSession.UserName},
		{"user_role", newSession.UserRole},
		{"access_token", newSession.AccessToken},
		{"refresh_token", newSession.RefreshToken},
		{"last_visit", newSession.LastVisit},
		{"mobile_key", newSession.MobileKey},
	})
	if errInsert != nil {
		return "", errInsert
	}
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (ar *authRepository) updateSession(access string, refresh string, entity *domain.UserEntity, authReq request.AuthenticateUserRequest) (string, error) {
	res := ar.collectionSessions.FindOneAndUpdate(
		utils.GetContext(),
		bson.D{
			{"user_id", entity.Id},
			{"user_name", entity.UserName},
			{"mobile_key", authReq.MobileKey},
		},
		bson.D{{
			"$set", bson.D{
				{"access_token", access},
				{"refresh_token", refresh},
				{"last_visit", utils.GetCurrentTime()},
			}}})

	var findSession domain.SessionEntity
	decodeResult := res.Decode(&findSession)
	if decodeResult != nil {
		return "", errors.Unwrap(fmt.Errorf("UpdateSession Decode ERROR"))
	}
	return findSession.Id, nil
}

func createSession(access string, refresh string, authReq request.AuthenticateUserRequest, entity *domain.UserEntity) domain.SessionEntity {
	return domain.SessionEntity{
		UserId:       entity.Id,
		UserName:     authReq.User.UserName,
		UserRole:     entity.UserRole,
		AccessToken:  access,
		RefreshToken: refresh,
		LastVisit:    utils.GetCurrentTime(),
		MobileKey:    authReq.MobileKey,
	}
}
