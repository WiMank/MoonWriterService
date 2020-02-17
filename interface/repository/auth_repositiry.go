package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WiMank/MoonWriterService/domain"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type authRepository struct {
	collectionUsers    *mongo.Collection
	collectionSessions *mongo.Collection
}

type AuthRepository interface {
	DecodeRequest(r *http.Request) request.AuthenticateUserRequest
	AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse
}

func NewAuthRepository(collectionUsers *mongo.Collection, collectionSessions *mongo.Collection) AuthRepository {
	return &authRepository{collectionUsers, collectionSessions}
}

func (ar *authRepository) DecodeRequest(r *http.Request) request.AuthenticateUserRequest {
	var requestUser request.AuthenticateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Error("Decode User response! ", err)
	}
	return requestUser
}

func (ar *authRepository) AuthenticateUser(authReq request.AuthenticateUserRequest) response.AppResponse {
	var localUserEntity domain.UserEntity
	userBson := bson.M{"user_name": authReq.User.UserName}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	if errFind := ar.collectionUsers.FindOne(ctx, userBson).Decode(&localUserEntity); errFind != nil {
		log.Errorf(fmt.Sprintf("AuthenticateUser could not find user: [%s]", authReq.User.UserName))
		//TODO: User not found
	}

	if localUserEntity.CheckUserCredentialsValid(authReq.User) {
		ar.checkSessionsCount(ctx, userBson)

	}

	unauthorized := response.UnauthorizedResponse{
		Message: "Unauthorized",
		Code:    http.StatusUnauthorized,
		Desc:    http.StatusText(http.StatusUnauthorized),
	}
	unauthorized.PrintLog(nil)

	return &unauthorized
}

func (ar *authRepository) checkSessionsCount(ctx context.Context, userBson bson.M) {
	count, errCount := ar.collectionSessions.CountDocuments(ctx, userBson)

	if count > 5 {
		ar.clearSessions(ctx, userBson)
	}

	if errCount != nil {
		log.Errorf("CountDocuments error: \n", errCount)
	}
}

func (ar *authRepository) clearSessions(ctx context.Context, userBson bson.M) {
	result, errDelete := ar.collectionSessions.DeleteMany(ctx, userBson)

	if errDelete != nil {
		log.Errorf("DeleteMany error: \n", errDelete)
	}

	if result != nil {
		log.Info("Delete result: ", result.DeletedCount)
	}
}

func createAccessToken(aur request.AuthenticateUserRequest) (string, int64) {
	tokenTime := GetCurrentTime()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    aur.User.UserName,
		"role":    aur.User.UserRole,
		"expired": tokenTime + 18e5,
	})

	tokenString, err := token.SignedString([]byte(domain.SecretKey))

	if err != nil {

	}

	return tokenString, tokenTime
}

func createRefreshToken(aur request.AuthenticateUserRequest) (string, int64) {
	tokenTime := GetCurrentTime()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"time":    tokenTime,
		"user":    aur.User.UserName,
		"expired": tokenTime + 2592e6,
	})

	tokenString, err := token.SignedString([]byte(domain.SecretKey))
	if err != nil {

	}

	return tokenString, tokenTime
}

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}
