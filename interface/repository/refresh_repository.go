package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WiMank/MoonWriterService/config"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
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
		return rr.responseCreator.CreateResponse(response.InvalidSession{}, "")
	}

	if localSession.RefreshToken == request.Refresh.RefreshToken {
		if rr.validateToken(request.Refresh.RefreshToken) {

		}
	}

	return rr.responseCreator.CreateResponse(response.InvalidSession{}, "")
}

func (rr *refreshRepository) findSession(request request.RefreshTokensRequest) (*domain.SessionEntity, error) {
	id, errHex := primitive.ObjectIDFromHex(request.Refresh.SessionId)
	if errHex != nil {
		return nil, errHex
	}

	var localSession domain.SessionEntity
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	findSessionErr := rr.collectionSessions.FindOne(ctx, bson.D{
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

	if token != nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["iss"], claims["user"], claims["exp"])
			return true
		} else {
			fmt.Println(err)
		}
	}
	return false
}

func (rr *refreshRepository) createAccessToken(entity domain.SessionEntity) (string, error) {
	tokenTime := getCurrentTime() + 36e2
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

func (rr *refreshRepository) createRefreshToken(entity domain.SessionEntity) (string, error) {
	tokenTime := getCurrentTime() + 2592e3
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
