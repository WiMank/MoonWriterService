package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WiMank/MoonWriterService/domain"
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
	DecodeUser(r *http.Request) domain.User
	EncodeUser(w http.ResponseWriter, response domain.UserResponse)
	AuthUser(user domain.User)
}

func NewAuthRepository(collectionUsers *mongo.Collection, collectionSessions *mongo.Collection) AuthRepository {
	return &authRepository{collectionUsers, collectionSessions}
}

func (ar *authRepository) DecodeUser(r *http.Request) domain.User {
	var requestUser domain.User
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		log.Error("Decode User response! ", err)
	}
	return requestUser
}

func (ar *authRepository) EncodeUser(w http.ResponseWriter, response domain.UserResponse) {
	w.WriteHeader(response.AppResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorf("Encode User response", err)
	}
}

func (ar *authRepository) AuthUser(user domain.User) {
	var localUser domain.User
	userBson := bson.M{"user_name": user.UserName}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	errFind := ar.collectionUsers.FindOne(ctx, userBson).Decode(&localUser)

	if errFind != nil {
		log.Errorf(fmt.Sprintf("AuthUser could not find user: [%s]", user.UserName))
	}

	if localUser.CheckUserCredentialsValid(user) {
		ar.checkSessionsCount(ctx, userBson)
	} else {
		//TODO: Unauthorized
	}

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
