package presenter

import (
	"encoding/json"
	"github.com/WiMank/MoonWriterService/interface/response"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type purchasePresenter struct {
}

type PurchasePresenter interface {
	PurchaseResponse(w http.ResponseWriter, appResponse response.AppResponse)
}

func NewPurchasePresenter() PurchasePresenter {
	return &purchasePresenter{}
}

func (pp *purchasePresenter) PurchaseResponse(w http.ResponseWriter, appResponse response.AppResponse) {
	w.WriteHeader(appResponse.GetStatusCode())
	err := json.NewEncoder(w).Encode(response.PurchaseResponse{AppResponse: appResponse})
	if err != nil {
		log.Errorf("PurchaseResponse error: \n", err)
	}
}
