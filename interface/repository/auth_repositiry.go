package repository

import (
	"github.com/WiMank/AlarmService/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type authRepository struct {
	collection *mongo.Collection
}

type AuthRepository interface {
	DecodeUser(r *http.Request) domain.User
}

func NewAuthRepository(collection *mongo.Collection) AuthRepository {
	return &authRepository{collection}
}

func (ar *authRepository) DecodeUser(r *http.Request) domain.User {
	return domain.User{}
}
