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
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type refreshRepository struct {
	collectionSessions *mongo.Collection
	responseCreator    response.AppResponseCreator
	validator          *validator.Validate
}

type RefreshRepository interface {
	DecodeRequest(r *http.Request) request.RefreshTokensRequest
	Refresh(request request.RefreshTokensRequest) response.AppResponse
}

func NewRefreshRepository(
	collectionSessions *mongo.Collection,
	responseCreator response.AppResponseCreator,
	validator *validator.Validate) RefreshRepository {
	return &refreshRepository{
		collectionSessions,
		responseCreator,
		validator,
	}
}

func (rr *refreshRepository) DecodeRequest(r *http.Request) request.RefreshTokensRequest {
	var refreshTokensRequest request.RefreshTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshTokensRequest); err != nil {
		log.Errorf("Decode RefreshTokensRequest error:\n", err)
	}
	return refreshTokensRequest
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
				return rr.responseCreator.CreateResponse(response.InvalidSession{}, request.Refresh.SessionId)
			}
		} else {
			return rr.responseCreator.CreateResponse(response.InvalidToken{}, "")
		}
	}

	return rr.responseCreator.CreateResponse(response.ValidateErrorResponse{}, "")
}

func (rr *refreshRepository) findSession(request request.RefreshTokensRequest) (*domain.SessionEntity, bool) {
	id, errHex := primitive.ObjectIDFromHex(request.Refresh.SessionId)

	if errHex != nil {
		return nil, false
	}

	var localSession domain.SessionEntity
	findSessionErr := rr.collectionSessions.FindOne(utils.GetContext(), bson.D{
		{"_id", id},
		{"mobile_key", request.Refresh.MobileKey},
		{"refresh_token", request.Refresh.RefreshToken},
	}).Decode(&localSession)

	if findSessionErr != nil {
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
	if entity != nil {

		id, errHex := primitive.ObjectIDFromHex(entity.Id)

		if errHex != nil {
			return false
		}

		_, errUpdate := rr.collectionSessions.UpdateOne(
			utils.GetContext(),
			bson.D{
				{"_id", id},
				{"refresh_token", entity.RefreshToken},
				{"mobile_key", entity.MobileKey},
			},
			bson.D{{
				"$set", bson.D{
					{"access_token", access},
					{"refresh_token", refresh},
					{"last_visit", utils.GetCurrentTime()},
				}}})

		if errUpdate != nil {
			return false
		}

		return true

	} else {
		return false
	}
}

func (rr *refreshRepository) createRefreshTokenResponse(entity *domain.SessionEntity, access string, refresh string) response.AppResponse {
	return rr.responseCreator.CreateResponse(response.TokenResponse{
		Message:      fmt.Sprintf("Tokens refreshed for [%s]", entity.UserName),
		SessionId:    entity.Id,
		AccessToken:  access,
		RefreshToken: refresh,
	}, "")
}
