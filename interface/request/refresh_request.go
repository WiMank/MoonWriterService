package request

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RefreshTokensRequest struct {
	Refresh struct {
		SessionId    int    `json:"session_id" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"required,gte=100,lte=200"`
		MobileKey    string `json:"mobile_key" validate:"required"`
	} `json:"refresh"`
}

func (rtr *RefreshTokensRequest) DecodeRefreshTokensRequest(r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&rtr); err != nil {
		log.Errorf("Decode RefreshTokensRequest error:\n", err)
	}
}

func (rtr *RefreshTokensRequest) ValidateRequest(validator *validator.Validate) bool {
	err := validator.Struct(rtr)
	if err != nil {
		log.Errorf("RefreshTokensRequest validate error: ", err)
		return false
	}
	return true
}
