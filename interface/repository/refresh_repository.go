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
	"github.com/spf13/cast"
	"time"
)

type refreshRepository struct {
	db              *sqlx.DB
	responseCreator response.AppResponseCreator
	validator       *validator.Validate
}

type RefreshRepository interface {
	Refresh(request request.RefreshTokensRequest) response.AppResponse
}

func NewRefreshRepository(
	db *sqlx.DB,
	responseCreator response.AppResponseCreator,
	validator *validator.Validate) RefreshRepository {
	return &refreshRepository{
		db,
		responseCreator,
		validator,
	}
}

func (rr *refreshRepository) Refresh(request request.RefreshTokensRequest) response.AppResponse {
	if request.ValidateRequest(rr.validator) {

		tokenValid := rr.validateToken(request.Refresh.RefreshToken)
		localSession, sessionExist := rr.findSession(request)
		accessToken, accessTokenCreated := rr.createAccessToken(localSession)
		refreshToken, refreshTokenCreated := rr.createRefreshToken(localSession)

		if tokenValid {
			if sessionExist {
				if accessTokenCreated && refreshTokenCreated {
					isRefreshed := rr.refreshSession(accessToken, refreshToken, localSession)
					if isRefreshed {
						return rr.createRefreshTokenResponse(localSession, accessToken, refreshToken)
					} else {
						return rr.responseCreator.CreateResponse(response.RefreshSessionErrorResponse{}, localSession.UserName)
					}
				} else {
					rr.responseCreator.CreateResponse(response.TokenErrorResponse{}, localSession.UserName)
				}
			} else {
				return rr.responseCreator.CreateResponse(response.InvalidSession{}, cast.ToString(request.Refresh.SessionId))
			}
		} else {
			return rr.responseCreator.CreateResponse(response.InvalidToken{}, "")
		}
	}

	return rr.responseCreator.CreateResponse(response.ValidateErrorResponse{}, "")
}

func (rr *refreshRepository) findSession(request request.RefreshTokensRequest) (*domain.SessionEntity, bool) {
	var localSession domain.SessionEntity
	existQuery := "SELECT * FROM sessions WHERE session_id=$1 AND mobile_key=$2 AND refresh_token=$3"
	err := rr.db.QueryRowx(existQuery,
		request.Refresh.SessionId,
		request.Refresh.MobileKey,
		request.Refresh.RefreshToken).StructScan(&localSession)
	if err != nil {
		return nil, false
	}

	return &localSession, true
}

func (rr *refreshRepository) validateToken(refreshToken string) bool {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return false
	}

	if token != nil && token.Valid {
		return true
	}

	return false
}

func (rr *refreshRepository) createAccessToken(entity *domain.SessionEntity) (string, bool) {
	if entity != nil {

		tokenTime := utils.GetAccessTokenTime()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":  "Moon Writer",
			"user": entity.UserId,
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

func (rr *refreshRepository) createRefreshToken(entity *domain.SessionEntity) (string, bool) {
	if entity != nil {

		tokenTime := utils.GetRefreshTokenTime()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":  "Moon Writer",
			"user": entity.UserId,
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

func (rr *refreshRepository) refreshSession(access string, refresh string, entity *domain.SessionEntity) bool {
	updateExec := "UPDATE sessions SET (access_token, refresh_token, last_visit) = ($1, $2, $3) " +
		"WHERE session_id=$4 AND refresh_token=$5 AND mobile_key=$6"
	_, err := rr.db.Exec(updateExec, access, refresh, time.Now(), entity.Id, entity.RefreshToken, entity.MobileKey)
	if err != nil {
		return false
	}

	return true
}

func (rr *refreshRepository) createRefreshTokenResponse(entity *domain.SessionEntity, access string, refresh string) response.AppResponse {
	return rr.responseCreator.CreateResponse(response.TokenResponse{
		Message:      fmt.Sprintf("Tokens refreshed for [%s]", entity.UserName),
		SessionId:    0,
		AccessToken:  access,
		RefreshToken: refresh,
	}, "")
}
