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
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type refreshRepository struct {
	collectionSessions *mongo.Collection
	responseCreator    response.AppResponseCreator
}

type RefreshRepository interface {
	DecodeRequest(r *http.Request) request.RefreshTokensRequest
	Refresh(request request.RefreshTokensRequest) response.AppResponse
}

func NewRefreshRepository(collectionSessions *mongo.Collection, responseCreator response.AppResponseCreator) RefreshRepository {
	return &refreshRepository{collectionSessions, responseCreator}
}

func (rr *refreshRepository) DecodeRequest(r *http.Request) request.RefreshTokensRequest {
	var refreshTokensRequest request.RefreshTokensRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshTokensRequest); err != nil {
		log.Errorf("Decode RefreshTokensRequest error:\n", err)
	}
	return refreshTokensRequest
}

func (rr *refreshRepository) Refresh(request request.RefreshTokensRequest) response.AppResponse {
	localSession, err := rr.findSession(request)
	if err != nil {
		return rr.responseCreator.CreateResponse(response.InvalidSession{}, request.Refresh.SessionId)
	}

	if rr.validateToken(request.Refresh.RefreshToken) {
		access, errAuthenticate := rr.createAccessToken(localSession)
		refresh, errAuthenticate := rr.createRefreshToken(localSession)
		if errAuthenticate != nil {
			return rr.responseCreator.CreateResponse(response.TokenErrorResponse{}, localSession.UserName)
		}

		errRefresh := rr.refreshSession(access, refresh, localSession)
		if errRefresh != nil {
			return rr.responseCreator.CreateResponse(response.RefreshSessionErrorResponse{}, localSession.UserName)
		}

		return rr.responseCreator.CreateResponse(response.TokenResponse{
			Message:      fmt.Sprintf("Tokens refreshed for [%s]", localSession.UserName),
			SessionId:    localSession.Id,
			AccessToken:  access,
			RefreshToken: refresh,
		}, "")
	}

	return rr.responseCreator.CreateResponse(response.InvalidToken{}, localSession.Id)
}

func (rr *refreshRepository) findSession(request request.RefreshTokensRequest) (*domain.SessionEntity, error) {
	id, errHex := primitive.ObjectIDFromHex(request.Refresh.SessionId)
	if errHex != nil {
		return nil, errHex
	}

	var localSession domain.SessionEntity
	findSessionErr := rr.collectionSessions.FindOne(utils.GetContext(), bson.D{
		{"_id", id},
		{"mobile_key", request.Refresh.MobileKey},
		{"refresh_token", request.Refresh.RefreshToken},
	}).Decode(&localSession)

	if findSessionErr != nil {
		return nil, findSessionErr
	}
	return &localSession, nil
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

func (rr *refreshRepository) createAccessToken(entity *domain.SessionEntity) (string, error) {
	tokenTime := utils.GetAccessTokenTime()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "Moon Writer",
		"user": entity.UserId,
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

func (rr *refreshRepository) createRefreshToken(entity *domain.SessionEntity) (string, error) {
	tokenTime := utils.GetRefreshTokenTime()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "Moon Writer",
		"user": entity.UserId,
		"exp":  tokenTime,
	})
	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Errorf("Refresh token signed error:\n", err)
		return "", err
	}
	return tokenString, nil
}

func (rr *refreshRepository) refreshSession(access string, refresh string, entity *domain.SessionEntity) error {
	id, errHex := primitive.ObjectIDFromHex(entity.Id)
	if errHex != nil {
		return errHex
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
		return errUpdate
	}
	return nil
}
