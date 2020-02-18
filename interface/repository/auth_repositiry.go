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
	var localUserEntity, errUser = ar.findUserEntity(authReq)
	var session, errSession = ar.findSession(authReq.MobileKey)
	if errUser != nil {
		return createUserFindResponse(errUser)
	}
	if errSession != nil {
		log.Error("FindSession error:\n", errSession)
	}

	if localUserEntity.CheckUserNameAndPass(authReq.User) {
		ar.checkSessionsCount(authReq)
		access := createAccessToken(authReq)
		refresh := createRefreshToken(authReq)
		if session.CheckMkExist(authReq.MobileKey) {
			ar.updateSession(refresh, access, authReq)
		} else {
			ar.insertSession(access, refresh, authReq)
		}
		return createTokenResponse(access, refresh, authReq)
	}
	return creteUnauthorizedResponse()
}

func (ar *authRepository) findUserEntity(authReq request.AuthenticateUserRequest) (*domain.UserEntity, error) {
	var localUserEntity domain.UserEntity
	userBson := bson.M{"user_name": authReq.User.UserName}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if errFind := ar.collectionUsers.FindOne(ctx, userBson).Decode(&localUserEntity); errFind != nil {
		return nil, errFind
	}
	return &localUserEntity, nil
}

func (ar *authRepository) findSession(mk string) (*domain.SessionEntity, error) {
	var session domain.SessionEntity
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	errMk := ar.collectionSessions.FindOne(ctx, bson.M{"mobile_key": mk}).Decode(&session)
	if errMk != nil {
		log.Error("Mobile Key decode error:\n", errMk)
		return nil, errMk
	}
	return &session, nil
}

func (ar *authRepository) checkSessionsCount(authReq request.AuthenticateUserRequest) {
	userBson := bson.M{"user_name": authReq.User.UserName}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	count, errCount := ar.collectionSessions.CountDocuments(ctx, userBson)
	if count > 5 {
		ar.clearSessions(ctx, userBson)
	}
	if errCount != nil {
		log.Errorf("CheckSessionsCount error: \n", errCount)
	}
}

func (ar *authRepository) clearSessions(ctx context.Context, userBson bson.M) {
	result, errDelete := ar.collectionSessions.DeleteMany(ctx, userBson)
	if errDelete != nil {
		log.Errorf("ClearSessions error: \n", errDelete)
	}
	if result != nil {
		log.Info("Delete result: ", result.DeletedCount)
	}
}

func createAccessToken(aur request.AuthenticateUserRequest) domain.Token {
	tokenTime := getCurrentTime() + 36e2
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    aur.User.UserName,
		"role":    aur.User.UserRole,
		"expired": tokenTime,
	})
	tokenString, err := token.SignedString([]byte(domain.SecretKey))
	if err != nil {

	}
	return domain.Token{
		Tok:     tokenString,
		Expired: tokenTime,
	}
}

func (ar *authRepository) insertSession(acc domain.Token, ref domain.Token, authReq request.AuthenticateUserRequest) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := ar.collectionSessions.InsertOne(ctx, createSession(acc, ref, authReq))
	if err != nil {
		log.Error("InsertSession error:\n", err)
	}
}

func createRefreshToken(aur request.AuthenticateUserRequest) domain.Token {
	tokenTime := getCurrentTime() + 2592e3
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"time":    tokenTime,
		"user":    aur.User.UserName,
		"expired": tokenTime,
	})
	tokenString, err := token.SignedString([]byte(domain.SecretKey))
	if err != nil {

	}
	return domain.Token{
		Tok:     tokenString,
		Expired: tokenTime,
	}
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}

func (ar *authRepository) updateSession(refresh domain.Token, access domain.Token, authReq request.AuthenticateUserRequest) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, errUpdate := ar.collectionSessions.UpdateOne(ctx,
		bson.D{{"user_name", authReq.User.UserName}, {"mobile_key", authReq.MobileKey}},
		bson.D{{"$set", bson.D{
			{"refresh_token", refresh.Tok},
			{"expires_in_r", refresh.Expired},
			{"access_token", access.Tok},
			{"expires_in_a", access.Expired},
			{"last_visit", getCurrentTime()},
		}}})

	if errUpdate != nil {
		log.Error("Refresh token ERROR: ", errUpdate)
	}
}

func createSession(acc domain.Token, ref domain.Token, authReq request.AuthenticateUserRequest) domain.SessionEntity {
	return domain.SessionEntity{
		UserName:     authReq.User.UserName,
		RefreshToken: ref.Tok,
		ExpiresInR:   ref.Expired,
		AccessToken:  acc.Tok,
		ExpiresInA:   acc.Expired,
		LastVisit:    getCurrentTime(),
		MobileKey:    authReq.MobileKey,
	}
}

func createUserFindResponse(err error) *response.UserFindResponse {
	s := response.UserFindResponse{
		Message: "User not found",
		Code:    http.StatusBadRequest,
		Desc:    http.StatusText(http.StatusBadRequest),
	}
	s.PrintLog(err)
	return &s
}

func createTokenResponse(acc domain.Token, ref domain.Token, authReq request.AuthenticateUserRequest) *response.TokenResponse {
	tokenResponse := response.TokenResponse{
		Message:      fmt.Sprintf("Tokens created for %s", authReq.User.UserName),
		Code:         http.StatusOK,
		Desc:         http.StatusText(http.StatusOK),
		RefreshToken: ref.Tok,
		ExpiresInR:   ref.Expired,
		AccessToken:  acc.Tok,
		ExpiresInA:   acc.Expired,
	}
	tokenResponse.PrintLog(nil)
	return &tokenResponse
}

func creteUnauthorizedResponse() *response.UnauthorizedResponse {
	unauthorized := response.UnauthorizedResponse{
		Message: "Unauthorized",
		Code:    http.StatusUnauthorized,
		Desc:    http.StatusText(http.StatusUnauthorized),
	}
	unauthorized.PrintLog(nil)
	return &unauthorized
}
