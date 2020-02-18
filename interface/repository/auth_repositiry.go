package repository

import (
	"context"
	"encoding/json"
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
	var errAuthenticate error
	localUserEntity, errAuthenticate := ar.findUserEntity(authReq)
	session, errAuthenticate := ar.findSession(authReq.MobileKey)

	if errAuthenticate != nil {
		return ar.responseCreator.CreateResponse(response.UserFindResponse{}, authReq.User.UserName)
	}

	if localUserEntity.CheckUserNameAndPass(authReq.User) {
		ar.checkSessionsCount(authReq)
		access, errAuthenticate := createAccessToken(authReq)
		refresh, errAuthenticate := createRefreshToken(authReq)
		if errAuthenticate != nil {
			return ar.responseCreator.CreateResponse(response.TokenErrorResponse{}, authReq.User.UserName)
		}

		if session.CheckMkExist(authReq.MobileKey) {
			ar.updateSession(access, refresh, authReq)
		} else {
			ar.insertSession(access, refresh, authReq)
		}

		return ar.responseCreator.CreateResponse(
			response.TokenResponse{
				AccessToken:  access.Tok,
				RefreshToken: refresh.Tok,
				ExpiresInA:   access.Expired,
				ExpiresInR:   refresh.Expired,
			},
			authReq.User.UserName,
		)
	}
	return ar.responseCreator.CreateResponse(response.UnauthorizedResponse{}, authReq.User.UserName)
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
		log.Errorf("ClearSessions error:\n", errDelete)
	}
	if result != nil {
		log.Info("Delete result: ", result.DeletedCount)
	}
}

func createAccessToken(aur request.AuthenticateUserRequest) (*domain.Token, error) {
	tokenTime := getCurrentTime() + 36e2
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    aur.User.UserName,
		"role":    aur.User.UserRole,
		"expired": tokenTime,
	})
	tokenString, err := token.SignedString([]byte(domain.SecretKey))
	if err != nil {
		log.Errorf("Access token signed error:\n", err)
		return nil, err
	}
	return &domain.Token{
		Tok:     tokenString,
		Expired: tokenTime,
	}, nil
}

func createRefreshToken(aur request.AuthenticateUserRequest) (*domain.Token, error) {
	tokenTime := getCurrentTime() + 2592e3
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"time":    tokenTime,
		"user":    aur.User.UserName,
		"expired": tokenTime,
	})
	tokenString, err := token.SignedString([]byte(domain.SecretKey))
	if err != nil {
		log.Errorf("Refresh token signed error:\n", err)
		return nil, err
	}
	return &domain.Token{
		Tok:     tokenString,
		Expired: tokenTime,
	}, nil
}

func (ar *authRepository) insertSession(access *domain.Token, refresh *domain.Token, authReq request.AuthenticateUserRequest) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := ar.collectionSessions.InsertOne(ctx, createSession(access, refresh, authReq))
	if err != nil {
		log.Errorf("InsertSession error:\n", err)
	}
}

func (ar *authRepository) updateSession(access *domain.Token, refresh *domain.Token, authReq request.AuthenticateUserRequest) {
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
		log.Errorf("UpdateSession error:\n", errUpdate)
	}
}

func createSession(access *domain.Token, refresh *domain.Token, authReq request.AuthenticateUserRequest) domain.SessionEntity {
	return domain.SessionEntity{
		UserName:     authReq.User.UserName,
		RefreshToken: refresh.Tok,
		ExpiresInR:   refresh.Expired,
		AccessToken:  access.Tok,
		ExpiresInA:   access.Expired,
		LastVisit:    getCurrentTime(),
		MobileKey:    authReq.MobileKey,
	}
}

func getCurrentTime() int64 {
	return time.Now().Unix()
}
