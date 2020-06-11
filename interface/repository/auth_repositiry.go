package repository

import (
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
)

type authRepository struct {
	db              *sqlx.DB
	responseCreator response.AppResponseCreator
	validator       *validator.Validate
}

type AuthRepository interface {
	AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse
}

func NewAuthRepository(db *sqlx.DB, responseCreator response.AppResponseCreator, validator *validator.Validate) AuthRepository {
	return &authRepository{db, responseCreator, validator}
}

func (ar *authRepository) AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse {
	if authReq.ValidateRequest(ar.validator) {
		go ar.checkSessionsCount(authReq)
		localUserEntity, userExist := ar.findUserEntity(authReq)
		passwordAndNameCorrect := localUserEntity.CheckUserNameAndPass(authReq.User)
		sessionExist := ar.checkSessionExist(authReq.MobileKey)
		accessToken, accessTokenCreated := createAccessToken(localUserEntity)
		refreshToken, refreshTokenCreated := createRefreshToken(localUserEntity)

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
						insertResult, sessionInserted := ar.insertSession(accessToken, refreshToken, authReq.MobileKey, localUserEntity)
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
	existQuery := "SELECT * FROM users WHERE user_name=$1"
	err := ar.db.QueryRowx(existQuery, authReq.User.UserName).StructScan(&localUserEntity)
	if err != nil {
		return nil, false
	}
	return &localUserEntity, true
}

func (ar *authRepository) checkSessionExist(mk string) bool {
	var exist bool
	existQuery := "SELECT EXISTS (SELECT FROM sessions WHERE mobile_key=$1)::boolean"
	err := ar.db.QueryRowx(existQuery, mk).Scan(&exist)
	if err != nil {
		return true
	}
	return exist
}

func (ar *authRepository) checkSessionsCount(authReq request.AuthenticateUserRequest) {
	var count int
	countQuery := "SELECT count(*) FROM sessions WHERE user_name=$1"
	err := ar.db.QueryRowx(countQuery, authReq.User.UserName).Scan(&count)

	if err != nil {
		log.Errorf("CheckSessionsCount error: \n", err)
	}

	if count > 5 {
		ar.clearSessions(authReq.User.UserName)
	}
}

func (ar *authRepository) clearSessions(userName string) {
	clearSessionsExec := "DELETE FROM sessions WHERE user_name =$1"
	result, err := ar.db.Exec(clearSessionsExec, userName)

	if err != nil {
		log.Errorf("ClearSessions error: ", err)
	}

	if result != nil {
		log.Info("Delete result: ", result)
	}
}

func createAccessToken(entity *domain.UserEntity) (string, bool) {
	if entity != nil {
		tokenTime := utils.GetAccessTokenTime()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":  "Moon Writer",
			"user": entity.UserName,
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

func createRefreshToken(entity *domain.UserEntity) (string, bool) {
	if entity != nil {
		tokenTime := utils.GetRefreshTokenTime()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":  "Moon Writer",
			"user": entity.UserName,
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

func (ar *authRepository) insertSession(
	access string,
	refresh string,
	mk string,
	entity *domain.UserEntity) (int, bool) {

	var id int
	insertQuery := "INSERT INTO sessions as s (user_name, user_role, access_token, refresh_token, mobile_key, last_visit) " +
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING s.session_id"
	err := ar.db.QueryRowx(insertQuery, entity.UserName, entity.UserRole, access, refresh, mk, utils.GetCurrentTime()).Scan(&id)

	if err != nil {
		log.Error("insertSession: ", err)
		return 0, false
	}

	return id, true
}

func (ar *authRepository) updateSession(
	access string,
	refresh string,
	entity *domain.UserEntity,
	authReq request.AuthenticateUserRequest,
) (int, bool) {

	var id int
	updateQuery := "UPDATE sessions as s SET (access_token, refresh_token, last_visit) = ($1, $2, $3)" +
		" WHERE user_name=$4 AND mobile_key=$5 RETURNING s.session_id"
	err := ar.db.QueryRowx(updateQuery, access, refresh, utils.GetCurrentTime(), entity.UserName, authReq.MobileKey).Scan(&id)

	if err != nil {
		return 0, false
	}

	return id, true
}

func (ar *authRepository) createUpdateTokenResponse(
	user *domain.UserEntity,
	sessionId int,
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
	sessionId int,
	accessToken string,
	refreshToken string) response.AppResponse {
	return ar.responseCreator.CreateResponse(response.TokenResponse{
		Message:      fmt.Sprintf("Tokens created for [%s]", user.UserName),
		SessionId:    sessionId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken}, user.UserName)
}
