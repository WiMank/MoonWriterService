package request

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PurchaseRegisterRequest struct {
	Purchase struct {
		AccessToken   string `json:"access_token" validate:"required,gte=100,lte=200"`
		PurchaseToken string `json:"purchase_token" validate:"required,gte=100,lte=200"`
		OrderId       string `json:"order_id" validate:"required,len=24"`
		PurchaseTime  int64  `json:"purchase_time" validate:"required"`
		Sku           string `json:"sku" validate:"required"`
	} `json:"purchase" validate:"required"`
}

type PurchaseVerificationRequest struct {
	Purchase struct {
		AccessToken string `json:"access_token" validate:"required,gte=100,lte=200"`
	} `json:"purchase" validate:"required"`
}

func (prr *PurchaseRegisterRequest) DecodePurchaseRegisterRequest(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&prr); err != nil {
		log.Errorf("DecodePurchaseRegisterRequest error:\n", err)
	}
}

func (pvr *PurchaseVerificationRequest) DecodeVerificationRequest(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&pvr); err != nil {
		log.Errorf("DecodeVerificationRequest error:\n", err)
	}
}

func (prr *PurchaseRegisterRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(prr)
	if err != nil {
		log.Errorf("PurchaseRegisterRequest validate error: ", err)
		return false
	}
	return true
}

func (pvr *PurchaseVerificationRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(pvr)
	if err != nil {
		log.Errorf("PurchaseVerificationRequest validate error: ", err)
		return false
	}
	return true
}
