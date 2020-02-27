package repository

import (
	"encoding/json"
	"fmt"
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/WiMank/MoonWriterService/interface/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type authRepository struct {
	db              *sqlx.DB
	responseCreator response.AppResponseCreator
	validator       *validator.Validate
}

type AuthRepository interface {
	DecodeRequest(r *http.Request) request.AuthenticateUserRequest
	AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse
}

func NewAuthRepository(db *sqlx.DB, responseCreator response.AppResponseCreator, validator *validator.Validate) AuthRepository {
	return &authRepository{db, responseCreator, validator}
}

func (ar *authRepository) DecodeRequest(r *http.Request) request.AuthenticateUserRequest {
	var requestUser request.AuthenticateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Errorf("Decode User response error:\n", err)
	}
	return requestUser
}

func (ar *authRepository) AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse {
	if authReq.ValidateRequest(ar.validator) {

		ar.checkSessionsCount(authReq)
		localUserEntity, userExist := ar.findUserEntity(authReq)
		passwordAndNameCorrect := localUserEntity.CheckUserNameAndPass(authReq.User)
		sessionExist := ar.checkSessionExist(authReq.MobileKey)
		accessToken, accessTokenCreated := ar.createAccessToken(localUserEntity)
		refreshToken, refreshTokenCreated := ar.createRefreshToken(localUserEntity)

		if userExist {
			if passwordAndNameCorrect {
				if accessTokenCreated && refreshTokenCreated {
					if sessionExist {
						updateResult, sessionUpdated := ar.updateSession(accessToken, refreshToken, localUserEntity, authReq)
						if sessionUpdated {
							return ar.createUpdateTokenResponse(localUserEntity, updateResult, accessToken, refreshToken)
						} else {
							return ar.responseCreator.CreateResponse(response.SessionUpdateFailedResponse{}, authReq.User.UserName)
						}
					} else {
						insertResult, sessionInserted := ar.insertSession(accessToken, refreshToken, authReq, localUserEntity)
						if sessionInserted {
							return ar.createNewTokenResponse(localUserEntity, insertResult, accessToken, refreshToken)
						} else {
							return ar.responseCreator.CreateResponse(response.SessionInsertFailedResponse{}, authReq.User.UserName)
						}
					}
				} else {
					return ar.responseCreator.CreateResponse(response.TokenErrorResponse{}, authReq.User.UserName)
				}
			} else {
				return ar.responseCreator.CreateResponse(response.UnauthorizedResponse{}, authReq.User.UserName)
			}
		} else {
			return ar.responseCreator.CreateResponse(response.UserFindResponse{}, authReq.User.UserName)
		}
	}

	return ar.responseCreator.CreateResponse(response.ValidateErrorResponse{}, "")
}

func (ar *authRepository) findUserEntity(authReq request.AuthenticateUserRequest) (*domain.UserEntity, bool) {
	var localUserEntity domain.UserEntity
	/*userBson := bson.D{{"user_name", authReq.User.UserName}, {"user_pass", authReq.User.UserPass}}

	if errFind := ar.collectionUsers.FindOne(utils.GetContext(), userBson).Decode(&localUserEntity); errFind != nil {
		return nil, false
	}*/

	return &localUserEntity, true
}

func (ar *authRepository) checkSessionExist(mk string) bool {
	/*count, err := ar.collectionSessions.CountDocuments(utils.GetContext(), bson.M{"mobile_key": mk})

	if err != nil {
		return false
	}

	if count != 1 {
		return false
	}
	*/
	return true
}

func (ar *authRepository) checkSessionsCount(authReq request.AuthenticateUserRequest) {
	/*	userBson := bson.M{"user_name": authReq.User.UserName}
		count, errCount := ar.collectionSessions.CountDocuments(utils.GetContext(), userBson)

		if count > 5 {
			ar.clearSessions(userBson)
		}

		if errCount != nil {
			log.Errorf("CheckSessionsCount error: \n", errCount)
		}*/
}

func (ar *authRepository) clearSessions(userBson bson.M) {
	/*result, errDelete := ar.collectionSessions.DeleteMany(utils.GetContext(), userBson)

	if errDelete != nil {
		log.Errorf("ClearSessions error: ", errDelete)
	}

	if result != nil {
		log.Info("Delete result: ", result.DeletedCount)
	}*/
}

func (ar *authRepository) createAccessToken(entity *domain.UserEntity) (string, bool) {
	if entity != nil {

		tokenTime := utils.GetAccessTokenTime()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":  "Moon Writer",
			"user": entity.Id,
			"role": entity.UserRole,
			"exp":  tokenTime,
		})
		tokenString, err := token.SignedString([]byte(config.SecretKey))

		if err != nil {
			return "", false
		}

		return tokenString, true

	} else {
		return "", false
	}
}

func (ar *authRepository) createRefreshToken(entity *domain.UserEntity) (string, bool) {
	if entity != nil {

		tokenTime := utils.GetRefreshTokenTime()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":  "Moon Writer",
			"user": entity.Id,
			"exp":  tokenTime,
		})
		tokenString, err := token.SignedString([]byte(config.SecretKey))

		if err != nil {
			return "", false
		}

		return tokenString, true

	} else {
		return "", false
	}
}

func (ar *authRepository) insertSession(access string, refresh string, authReq request.AuthenticateUserRequest, entity *domain.UserEntity) (string, bool) {
	/*	newSession := createSession(access, refresh, authReq, entity)
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
			return "", false
		}

		return insertResult.InsertedID.(primitive.ObjectID).Hex(), true*/
	return "", false
}

func (ar *authRepository) updateSession(
	access string,
	refresh string,
	entity *domain.UserEntity,
	authReq request.AuthenticateUserRequest,
) (string, bool) {
	/*res := ar.collectionSessions.FindOneAndUpdate(
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
		return "", false
	}

	return findSession.Id, true*/
	return "", false
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

func (ar *authRepository) createUpdateTokenResponse(
	user *domain.UserEntity,
	sessionId string,
	accessToken string,
	refreshToken string) response.AppResponse {
	return ar.responseCreator.CreateResponse(response.TokenResponse{
		Message:      fmt.Sprintf("Tokens updated for [%s]", user.UserName),
		SessionId:    sessionId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken}, user.UserName)
}

func (ar *authRepository) createNewTokenResponse(
	user *domain.UserEntity,
	sessionId string,
	accessToken string,
	refreshToken string) response.AppResponse {
	return ar.responseCreator.CreateResponse(response.TokenResponse{
		Message:      fmt.Sprintf("Tokens created for [%s]", user.UserName),
		SessionId:    sessionId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken}, user.UserName)
}
