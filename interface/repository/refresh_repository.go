package repository

import (
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type refreshRepository struct {
	collectionSessions *mongo.Collection
	responseCreator    response.AppResponseCreator
}

type RefreshRepository interface {
	DecodeRequest(r *http.Request) request.AuthenticateUserRequest
	Refresh() response.AppResponse
}

func NewRefreshRepository(collectionUsers *mongo.Collection, collectionSessions *mongo.Collection, responseCreator response.AppResponseCreator) RefreshRepository {
	return &refreshRepository{collectionSessions, responseCreator}
}

func (rr *refreshRepository) DecodeRequest(r *http.Request) request.AuthenticateUserRequest {
	return request.AuthenticateUserRequest{}
}

func (rr *refreshRepository) Refresh() response.AppResponse {
	return &response.UnauthorizedResponse{}
}
