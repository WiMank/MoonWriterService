package repository

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/interface/request"
	"github.com/WiMank/MoonWriterService/interface/response"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type purchaseRepository struct {
	collectionSessions *mongo.Collection
	responseCreator    response.AppResponseCreator
}

type PurchaseRepository interface {
	DecodeRequest(r *http.Request) request.PurchaseRequest
	PurchaseVerification(request request.PurchaseRequest) response.AppResponse
}

func NewPurchaseRepository(collectionSessions *mongo.Collection, responseCreator response.AppResponseCreator) PurchaseRepository {
	return &purchaseRepository{collectionSessions, responseCreator}
}

func (pr *purchaseRepository) DecodeRequest(r *http.Request) request.PurchaseRequest {
	var refreshTokensRequest request.PurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshTokensRequest); err != nil {
		log.Errorf("Decode PurchaseVerification error:\n", err)
	}
	return refreshTokensRequest
}

func (pr *purchaseRepository) PurchaseVerification(request request.PurchaseRequest) response.AppResponse {
	return nil
}
